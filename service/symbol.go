package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/grpclog"

	"github.com/vectorman1/analysis/analysis-api/common"

	db2 "github.com/vectorman1/analysis/analysis-api/model/db"

	"github.com/vectorman1/analysis/analysis-api/generated/worker_symbol_service"

	"github.com/jackc/pgx"

	"github.com/jackc/pgx/pgtype"
	"github.com/vectorman1/analysis/analysis-api/db"
	"github.com/vectorman1/analysis/analysis-api/generated/proto_models"
	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
)

type symbolsService interface {
	// repo methods
	GetPaged(ctx *context.Context, req *symbol_service.ReadPagedSymbolRequest) (*[]*proto_models.Symbol, uint, error)
	Details(ctx *context.Context, req *symbol_service.SymbolDetailsRequest) (*symbol_service.SymbolDetailsResponse, error)
	InsertBulk(ctx *context.Context, symbols []*proto_models.Symbol) (bool, error)
	// service methods
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

func (s *SymbolsService) Recalculate(ctx *context.Context) (*symbol_service.RecalculateSymbolResponse, error) {
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
	ctx *context.Context) (*symbol_service.RecalculateSymbolResponse, error) {

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

	fmt.Println("starting save of update")

	timeoutContext, c := context.WithTimeout(*ctx, 10*time.Second)
	defer c()

	tx, err := s.symbolsRepository.BeginTx(&timeoutContext, &pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	fmt.Println("getting symbols from protos")
	// create new symbols
	createEntities, err := s.symbolDataToEntity(&createSymbols)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}

	fmt.Println("inserting new symbols")
	_, err = s.symbolsRepository.InsertBulk(tx, &timeoutContext, createEntities)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}

	fmt.Println("deleting old symbols")
	// delete entities
	deleteEntities, err := s.symbolDataToEntity(&deleteSymbols)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}
	_, err = s.symbolsRepository.DeleteBulk(tx, &timeoutContext, deleteEntities)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}

	fmt.Println("updating old symbols")
	// update entities
	updateEntities, err := s.symbolDataToEntity(&updateSymbols)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}
	_, err = s.symbolsRepository.UpdateBulk(tx, &timeoutContext, updateEntities)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}

	fmt.Println("saving")
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

func (s *SymbolsService) symbolDataToEntity(in *[]*proto_models.Symbol) ([]*db2.Symbol, error) {
	var result []*db2.Symbol
	var wg sync.WaitGroup

	wg.Add(len(*in))
	for _, sym := range *in {
		e := db2.Symbol{}

		// get or create the currency of the symbol
		curr, err := s.currencyRepository.GetOrCreate(sym.Currency.Code)
		if err != nil {
			return nil, err
		}

		e.Uuid = pgtype.UUID{}
		err = e.Uuid.Set(sym.Uuid)
		if err != nil {
			return nil, err
		}

		e.CurrencyID = curr.ID
		e.Isin = sym.Isin
		e.Identifier = sym.Identifier
		e.Name = sym.Name
		var moq pgtype.Float4
		_ = moq.Set(sym.MinimumOrderQuantity)
		e.MinimumOrderQuantity = moq
		e.MarketName = sym.MarketName
		e.MarketHoursGmt = sym.MarketHoursGmt

		e.CreatedAt = pgtype.Timestamptz{Time: time.Now(), Status: pgtype.Present}
		e.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Status: pgtype.Present}
		e.DeletedAt = pgtype.Timestamptz{Status: pgtype.Null}

		result = append(result, &e)
	}

	return result, nil
}

func (s *SymbolsService) generateRecalculationResult(newSymbols []*proto_models.Symbol, oldSymbols []*proto_models.Symbol) []*worker_symbol_service.RecalculateSymbolsResponse {
	unique := make(map[string]*worker_symbol_service.RecalculateSymbolsResponse)

	// generate create, update and ignore responses
	for _, newSym := range newSymbols {
		if ok, oldSym := common.ContainsSymbol(newSym.Uuid, oldSymbols); !ok {
			if unique[newSym.Uuid] == nil {
				unique[newSym.Uuid] =
					&worker_symbol_service.RecalculateSymbolsResponse{
						Type:   worker_symbol_service.RecalculateSymbolsResponse_CREATE,
						Symbol: newSym,
					}
				continue
			} else {
				log.Println("collision on create: ", newSym)
				log.Println("existing: ", unique[newSym.Uuid].Type, unique[newSym.Uuid])
			}
		} else {
			// check if any fields from the new symbol are different from the old
			shouldUpdate := false
			if oldSym.Name != newSym.Name {
				shouldUpdate = true
			} else if oldSym.MarketName != newSym.MarketName {
				shouldUpdate = true
			} else if oldSym.MarketHoursGmt != newSym.MarketHoursGmt {
				shouldUpdate = true
			}

			// if any fields are updated, send and update response, otherwise, send it back and ignore it
			if shouldUpdate {
				if unique[newSym.Uuid] == nil {
					unique[newSym.Uuid] =
						&worker_symbol_service.RecalculateSymbolsResponse{
							Type:   worker_symbol_service.RecalculateSymbolsResponse_UPDATE,
							Symbol: newSym,
						}
					continue
				} else {
					grpclog.Infoln("collision on update: ", newSym)
					grpclog.Infoln("existing: ", unique[newSym.Uuid].Type, unique[newSym.Uuid])
				}
			} else {
				if unique[newSym.Uuid] == nil {
					unique[newSym.Uuid] =
						&worker_symbol_service.RecalculateSymbolsResponse{
							Type:   worker_symbol_service.RecalculateSymbolsResponse_IGNORE,
							Symbol: newSym,
						}
					continue
				} else {
					grpclog.Infoln("collision on ignore: ", newSym)
					grpclog.Infoln("existing: ", unique[newSym.Uuid].Type, unique[newSym.Uuid])
				}
			}
		}
	}

	// generate delete responses
	for _, oldSym := range oldSymbols {
		if ok, _ := common.ContainsSymbol(oldSym.Uuid, newSymbols); !ok {
			if unique[oldSym.Uuid] == nil {
				unique[oldSym.Uuid] = &worker_symbol_service.RecalculateSymbolsResponse{
					Type:   worker_symbol_service.RecalculateSymbolsResponse_DELETE,
					Symbol: oldSym,
				}
			} else {
				grpclog.Infoln("collision: ", oldSym)
				grpclog.Infoln("existing: ", unique[oldSym.Uuid].Type, unique[oldSym.Uuid])
			}
		}
	}

	var result []*worker_symbol_service.RecalculateSymbolsResponse
	for _, v := range unique {
		result = append(result, v)
	}

	return result
}
