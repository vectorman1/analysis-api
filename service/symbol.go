package service

import (
	"context"
	"time"

	"google.golang.org/grpc/grpclog"

	"github.com/vectorman1/analysis/analysis-api/common"

	"github.com/vectorman1/analysis/analysis-api/generated/worker_symbol_service"

	"github.com/vectorman1/analysis/analysis-api/third_party/alpha_vantage"
	"github.com/vectorman1/analysis/analysis-api/third_party/trading_212"

	"github.com/vectorman1/analysis/analysis-api/model/db/entities"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jackc/pgx"

	"github.com/vectorman1/analysis/analysis-api/db"
	"github.com/vectorman1/analysis/analysis-api/generated/proto_models"
	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
)

type SymbolServiceContract interface {
	// repo methods
	Get(ctx context.Context, uuid string) (*proto_models.Symbol, error)
	GetPaged(ctx context.Context, req *symbol_service.GetPagedRequest) (*[]*proto_models.Symbol, uint, error)
	Overview(ctx context.Context, req *symbol_service.SymbolOverviewRequest) (*symbol_service.SymbolOverview, error)

	// service methods
	UpdateAll(ctx context.Context) (*symbol_service.RecalculateSymbolResponse, error)
	processRecalculationResponse(
		input []*worker_symbol_service.RecalculateSymbolsResponse,
		ctx context.Context) (*symbol_service.RecalculateSymbolResponse, error)
	symbolDataToEntity(in *[]*proto_models.Symbol) ([]*entities.Symbol, error)
	filterUnusableSymbols(symbols *[]*proto_models.Symbol) *[]*proto_models.Symbol
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

func (s *SymbolService) GetPaged(ctx context.Context, req *symbol_service.GetPagedRequest) (*[]*proto_models.Symbol, uint, error) {
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

func (s *SymbolService) UpdateAll(ctx context.Context) (*symbol_service.RecalculateSymbolResponse, error) {
	oldSymbols, _, err := s.GetPaged(ctx, &symbol_service.GetPagedRequest{
		Filter: &proto_models.PagedFilter{
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

	externalSymbols = s.filterUnusableSymbols(externalSymbols)
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
	var result []*worker_symbol_service.RecalculateSymbolsResponse
	unique := make(map[string]*worker_symbol_service.RecalculateSymbolsResponse)
	output := make(chan *worker_symbol_service.RecalculateSymbolsResponse)

	go func() {
		for _, newSym := range newSymbols {
			if ok, oldSym := common.ContainsSymbol(newSym.Uuid, oldSymbols); !ok {
				output <- &worker_symbol_service.RecalculateSymbolsResponse{
					Type:   worker_symbol_service.RecalculateSymbolsResponse_CREATE,
					Symbol: newSym,
				}
				continue
			} else {
				shouldUpdate := false
				if oldSym.Name != newSym.Name {
					shouldUpdate = true
				} else if oldSym.MarketHoursGmt != newSym.MarketHoursGmt {
					shouldUpdate = true
				}
				// Identifier, ISIN and Market Name are not checked, as they are used in the uuids of the symbols
				// if any fields are updated, send an update response
				if shouldUpdate {
					output <- &worker_symbol_service.RecalculateSymbolsResponse{
						Type:   worker_symbol_service.RecalculateSymbolsResponse_UPDATE,
						Symbol: newSym,
					}
					continue
				} else {
					output <- &worker_symbol_service.RecalculateSymbolsResponse{
						Type:   worker_symbol_service.RecalculateSymbolsResponse_IGNORE,
						Symbol: newSym,
					}
				}
			}
		}

		for _, sym := range oldSymbols {
			if ok, _ := common.ContainsSymbol(sym.Uuid, newSymbols); !ok {
				output <- &worker_symbol_service.RecalculateSymbolsResponse{
					Type:   worker_symbol_service.RecalculateSymbolsResponse_DELETE,
					Symbol: sym,
				}
			}
		}

		close(output)
	}()

	for res := range output {
		if existing := unique[res.Symbol.Uuid]; existing == nil {
			unique[res.Symbol.Uuid] = res
		} else {
			grpclog.Infof(
				"symbol update collision - new type: %d - existing type: %d (new, existing): %s, %s - %s, %s",
				res.Type, existing.Type,
				res.Symbol.Name, existing.Symbol.Name,
				res.Symbol.MarketHoursGmt, existing.Symbol.MarketHoursGmt)
		}
	}

	for _, res := range unique {
		result = append(result, res)
	}

	return result
}

func (s *SymbolService) filterUnusableSymbols(symbols *[]*proto_models.Symbol) *[]*proto_models.Symbol {
	var res []*proto_models.Symbol

	for _, sym := range *symbols {
		switch sym.MarketName {
		case common.MarketNASDAQ:
			res = append(res, sym)
		case common.MarketNYSE:
			res = append(res, sym)
		case common.MarketNonISANYSE:
			res = append(res, sym)
		case common.MarketNonISAOTCMarkets:
			res = append(res, sym)
		default:
			continue
		}
	}

	return &res
}
