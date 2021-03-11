package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/grpclog"

	"github.com/vectorman1/analysis/analysis-api/common"

	db2 "github.com/vectorman1/analysis/analysis-api/model/db"

	"github.com/vectorman1/analysis/analysis-api/generated/worker_symbol_service"

	"github.com/jackc/pgx"

	"github.com/vectorman1/analysis/analysis-api/db"
	"github.com/vectorman1/analysis/analysis-api/generated/proto_models"
	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
)

type symbolsService interface {
	// repo methods
	GetPaged(ctx *context.Context, req *symbol_service.ReadPagedSymbolRequest) (*[]*proto_models.Symbol, uint, error)
	Details(ctx *context.Context, req *symbol_service.SymbolDetailsRequest) (*symbol_service.SymbolDetailsResponse, error)

	// service methods
	Recalculate(ctx *context.Context) (*symbol_service.RecalculateSymbolResponse, error)
	processRecalculationResponse(
		input []*worker_symbol_service.RecalculateSymbolsResponse,
		ctx *context.Context) (*symbol_service.RecalculateSymbolResponse, error)
	symbolDataToEntity(in *[]*proto_models.Symbol) ([]*db2.Symbol, error)
}

type SymbolsService struct {
	symbolsService
	symbolsRepository        *db.SymbolRepository
	symbolOverviewRepository *db.SymbolOverviewRepository
	alphaVantageService      *AlphaVantageService
	externalSymbolService    *ExternalSymbolService
}

func NewSymbolsService(
	symbolsRepository *db.SymbolRepository,
	symbolOverviewRepository *db.SymbolOverviewRepository,
	alphaVantageService *AlphaVantageService,
	externalSymbolService *ExternalSymbolService) *SymbolsService {
	return &SymbolsService{
		symbolsRepository:        symbolsRepository,
		symbolOverviewRepository: symbolOverviewRepository,
		alphaVantageService:      alphaVantageService,
		externalSymbolService:    externalSymbolService,
	}
}

func (s *SymbolsService) GetPaged(ctx *context.Context, req *symbol_service.ReadPagedSymbolRequest) (*[]*proto_models.Symbol, uint, error) {
	var res []*proto_models.Symbol
	syms, totalItemsCount, err := s.symbolsRepository.GetPaged(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	for _, sym := range *syms {
		res = append(res, sym.ToProtoObject())
	}

	return &res, totalItemsCount, nil
}

func (s *SymbolsService) Details(ctx *context.Context, req *symbol_service.SymbolDetailsRequest) (*symbol_service.SymbolDetailsResponse, error) {
	symbol, err := s.symbolsRepository.GetByUuid(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	overview, err := s.symbolOverviewRepository.GetBySymbolUuid(ctx, req.Uuid)
	if err != nil {
		extOverview, err := s.alphaVantageService.GetOrUpdateSymbolOverview(symbol.Identifier)
		if err != nil {
			return nil, err
		}

		_, err = s.symbolOverviewRepository.Insert(ctx, extOverview.ToEntity(req.Uuid))
		if err != nil {
			return nil, err
		}

		overview, err = s.symbolOverviewRepository.GetBySymbolUuid(ctx, req.Uuid)
		if err != nil {
			return nil, err
		}
	}

	return &symbol_service.SymbolDetailsResponse{
		Symbol:  symbol.ToProtoObject(),
		Details: overview.ToProtoObject(),
	}, nil
}

func (s *SymbolsService) Recalculate(ctx context.Context) (*symbol_service.RecalculateSymbolResponse, error) {
	oldSymbols, _, err := s.GetPaged(ctx, &symbol_service.ReadPagedSymbolRequest{
		Filter: &symbol_service.SymbolFilter{
			PageSize:   100000,
			PageNumber: 1,
			Order:      "identifier",
			Ascending:  true,
		},
	})
	if err != nil {
		return nil, err
	}

	externalSymbols, err := s.externalSymbolService.GetLatest(ctx)
	if err != nil {
		return nil, err
	}

	result := s.generateRecalculationResult(*externalSymbols, *oldSymbols)
	response, err := s.processRecalculationResponse(result, ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *SymbolsService) processRecalculationResponse(
	input []*worker_symbol_service.RecalculateSymbolsResponse,
	ctx context.Context) (*symbol_service.RecalculateSymbolResponse, error) {

	var createSymbols []*proto_models.Symbol
	var updateSymbols []*proto_models.Symbol
	var deleteSymbols []*proto_models.Symbol

	itemsIgnored := int64(0)
	for _, res := range input {
		switch res.Type {
		case worker_symbol_service.RecalculateSymbolsResponse_CREATE:
			createSymbols = append(createSymbols, res.Symbol)
		case worker_symbol_service.RecalculateSymbolsResponse_UPDATE:
			updateSymbols = append(updateSymbols, res.Symbol)
		case worker_symbol_service.RecalculateSymbolsResponse_DELETE:
			deleteSymbols = append(deleteSymbols, res.Symbol)
		case worker_symbol_service.RecalculateSymbolsResponse_IGNORE:
			itemsIgnored++
		}
	}

	timeoutContext, c := context.WithTimeout(ctx, 10*time.Second)
	defer c()

	tx, err := s.symbolsRepository.BeginTx(&timeoutContext, &pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	// create new symbols
	createEntities, err := s.symbolDataToEntity(&createSymbols)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}

	_, err = s.symbolsRepository.InsertBulk(tx, timeoutContext, createEntities)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}

	// delete entities
	deleteEntities, err := s.symbolDataToEntity(&deleteSymbols)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}
	_, err = s.symbolsRepository.DeleteBulk(tx, timeoutContext, deleteEntities)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}

	// update entities
	updateEntities, err := s.symbolDataToEntity(&updateSymbols)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}
	_, err = s.symbolsRepository.UpdateBulk(tx, timeoutContext, updateEntities)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}

	err = tx.CommitEx(timeoutContext)
	if err != nil {
		return nil, err
	}

	return &symbol_service.RecalculateSymbolResponse{
		ItemsCreated: int64(len(createSymbols)),
		ItemsUpdated: int64(len(updateSymbols)),
		ItemsDeleted: int64(len(deleteSymbols)),
		ItemsIgnored: itemsIgnored,
		TotalItems:   int64(len(input)),
	}, nil
}

func (s *SymbolsService) symbolDataToEntity(in *[]*proto_models.Symbol) ([]*entities.Symbol, error) {
	var result []*entities.Symbol

	for _, sym := range *in {
		result = append(result, entities.Symbol{}.FromProtoObject(sym))
	}

	return result, nil
}

func (s *SymbolsService) generateRecalculationResult(newSymbols []*proto_models.Symbol, oldSymbols []*proto_models.Symbol) []*worker_symbol_service.RecalculateSymbolsResponse {
	var wg sync.WaitGroup
	unique := sync.Map{}
	// generate create, update and ignore responses
	for _, newSym := range newSymbols {
		go func(wg *sync.WaitGroup, s *proto_models.Symbol) {
			wg.Add(1)
			defer wg.Done()

			mapValue, exists := unique.Load(s.Uuid)
			if !exists {
				if ok, oldSym := common.ContainsSymbol(s.Uuid, oldSymbols); !ok {
					if !ok {
						wg.Add(1)
						defer wg.Done()
						unique.Store(s.Uuid,
							&worker_symbol_service.RecalculateSymbolsResponse{
								Type:   worker_symbol_service.RecalculateSymbolsResponse_CREATE,
								Symbol: s,
							})
						return
					}
				} else {
					// check if any fields from the new symbol are different from the old
					shouldUpdate := false
					if oldSym.Name != s.Name {
						shouldUpdate = true
					} else if oldSym.MarketHoursGmt != s.MarketHoursGmt {
						shouldUpdate = true
					}

					// Identifier, ISIN and Market Name are not checked, as they are used in the uuids of the symbols
					// if any fields are updated, send and update response, otherwise, send it back and ignore it
					if shouldUpdate {
						unique.Store(s.Uuid,
							&worker_symbol_service.RecalculateSymbolsResponse{
								Type:   worker_symbol_service.RecalculateSymbolsResponse_UPDATE,
								Symbol: s,
							})
						return
					} else {
						unique.Store(s.Uuid,
							&worker_symbol_service.RecalculateSymbolsResponse{
								Type:   worker_symbol_service.RecalculateSymbolsResponse_IGNORE,
								Symbol: s,
							})
						return
					}
				}
			} else {
				grpclog.Infoln("collision while checking for CREATE , UPDATE OR IGNORE - existing:", mapValue)
				grpclog.Infoln("new: ", s)
			}
		}(&wg, newSym)
	}

	// generate delete responses
	for _, oldSym := range oldSymbols {
		go func(wg *sync.WaitGroup, s *proto_models.Symbol) {
			wg.Add(1)
			defer wg.Done()

			mapValue, exists := unique.Load(s.Uuid)
			if !exists {
				if ok, _ := common.ContainsSymbol(s.Uuid, newSymbols); !ok {
					unique.Store(s.Uuid,
						&worker_symbol_service.RecalculateSymbolsResponse{
							Type:   worker_symbol_service.RecalculateSymbolsResponse_DELETE,
							Symbol: s,
						})
					return
				}
			} else {
				grpclog.Infoln("collision while checking for DELETE - existing:", mapValue)
				grpclog.Infoln("new: ", s)
			}
		}(&wg, oldSym)
	}

	wg.Wait()

	var result []*worker_symbol_service.RecalculateSymbolsResponse
	unique.Range(func(k, v interface{}) bool {
		result = append(result, v.(*worker_symbol_service.RecalculateSymbolsResponse))
		return true
	})

	return result
}
