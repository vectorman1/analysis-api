package service

import (
	"context"
	"sync"
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/worker_symbol_service"

	"github.com/vectorman1/analysis/analysis-api/third_party/alpha_vantage"
	"github.com/vectorman1/analysis/analysis-api/third_party/trading_212"

	"github.com/vectorman1/analysis/analysis-api/model/db/entities"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/grpclog"

	"github.com/vectorman1/analysis/analysis-api/common"

	"github.com/jackc/pgx"

	"github.com/vectorman1/analysis/analysis-api/db"
	"github.com/vectorman1/analysis/analysis-api/generated/proto_models"
	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
)

type symbolService interface {
	// repo methods
	Get(ctx context.Context, uuid string) (*proto_models.Symbol, error)
	GetPaged(ctx context.Context, req *symbol_service.ReadPagedRequest) (*[]*proto_models.Symbol, uint, error)
	Overview(ctx context.Context, req *symbol_service.SymbolOverviewRequest) (*symbol_service.SymbolOverview, error)

	// service methods
	Recalculate(ctx context.Context) (*symbol_service.RecalculateSymbolResponse, error)
	processRecalculationResponse(
		input []*worker_symbol_service.RecalculateSymbolsResponse,
		ctx context.Context) (*symbol_service.RecalculateSymbolResponse, error)
	symbolDataToEntity(in *[]*proto_models.Symbol) ([]*entities.Symbol, error)
}

type SymbolService struct {
	symbolRepository         *db.SymbolRepository
	symbolOverviewRepository *db.SymbolOverviewRepository
	alphaVantageService      *alpha_vantage.AlphaVantageService
	externalSymbolService    *trading_212.ExternalSymbolService
}

func (s *SymbolService) Get(ctx context.Context, uuid string) (*proto_models.Symbol, error) {
	sym, err := s.symbolRepository.GetByUuid(ctx, uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid symbol uuid")
	}

	return sym.ToProto(), nil
}

func NewSymbolService(
	symbolsRepository *db.SymbolRepository,
	symbolOverviewRepository *db.SymbolOverviewRepository,
	alphaVantageService *alpha_vantage.AlphaVantageService,
	externalSymbolService *trading_212.ExternalSymbolService) *SymbolService {
	return &SymbolService{
		symbolRepository:         symbolsRepository,
		symbolOverviewRepository: symbolOverviewRepository,
		alphaVantageService:      alphaVantageService,
		externalSymbolService:    externalSymbolService,
	}
}

func (s *SymbolService) GetPaged(ctx context.Context, req *symbol_service.ReadPagedRequest) (*[]*proto_models.Symbol, uint, error) {
	if req.Filter == nil {
		return nil, 0, status.Errorf(codes.InvalidArgument, "provide filter")
	}
	if req.Filter.Order == "" {
		return nil, 0, status.Error(codes.InvalidArgument, "provide order argument")
	}

	var res []*proto_models.Symbol
	syms, totalItemsCount, err := s.symbolRepository.GetPaged(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	for _, sym := range *syms {
		res = append(res, sym.ToProto())
	}

	return &res, totalItemsCount, nil
}

func (s *SymbolService) Overview(ctx context.Context, req *symbol_service.SymbolOverviewRequest) (*symbol_service.SymbolOverview, error) {
	userInfo := ctx.Value("user_info")
	if userInfo == nil {
		return nil, status.Error(codes.Unauthenticated, "provide user token")
	}

	symbol, err := s.symbolRepository.GetByUuid(ctx, req.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid symbol uuid")
	}

	overview, err := s.symbolOverviewRepository.GetBySymbolUuid(ctx, req.Uuid)
	if err != nil || overview.ShouldUpdate() {
		extOverview, err := s.alphaVantageService.GetSymbolOverview(symbol.Identifier)
		if err != nil {
			return nil, err
		}

		_, err = s.symbolOverviewRepository.Insert(ctx, extOverview.ToEntity(req.Uuid))
		if err != nil {
			return nil, err
		}

		newOverview, err := s.symbolOverviewRepository.GetBySymbolUuid(ctx, req.Uuid)
		if err != nil {
			return nil, err
		}
		overview = newOverview
	}

	return overview.ToProto(), nil
}

func (s *SymbolService) Recalculate(ctx context.Context) (*symbol_service.RecalculateSymbolResponse, error) {
	oldSymbols, _, err := s.GetPaged(ctx, &symbol_service.ReadPagedRequest{
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

func (s *SymbolService) processRecalculationResponse(
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

	tx, err := s.symbolRepository.BeginTx(&timeoutContext, &pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	// create new symbols
	createEntities, err := s.symbolDataToEntity(&createSymbols)
	if err != nil {
		tx.RollbackEx(timeoutContext)
		return nil, err
	}

	_, err = s.symbolRepository.InsertBulk(tx, timeoutContext, createEntities)
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
	_, err = s.symbolRepository.DeleteBulk(tx, timeoutContext, deleteEntities)
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
	_, err = s.symbolRepository.UpdateBulk(tx, timeoutContext, updateEntities)
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

func (s *SymbolService) symbolDataToEntity(in *[]*proto_models.Symbol) ([]*entities.Symbol, error) {
	var result []*entities.Symbol

	for _, sym := range *in {
		result = append(result, entities.Symbol{}.FromProtoObject(sym))
	}

	return result, nil
}

func (s *SymbolService) generateRecalculationResult(newSymbols []*proto_models.Symbol, oldSymbols []*proto_models.Symbol) []*worker_symbol_service.RecalculateSymbolsResponse {
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
					unique.Store(s.Uuid,
						&worker_symbol_service.RecalculateSymbolsResponse{
							Type:   worker_symbol_service.RecalculateSymbolsResponse_CREATE,
							Symbol: s,
						})
					return
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
