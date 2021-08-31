package service

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/model"
	"github.com/vectorman1/analysis/analysis-api/domain/instrument/repo"
	"github.com/vectorman1/analysis/analysis-api/domain/instrument/third_party"

	validationErrors "github.com/vectorman1/analysis/analysis-api/common/errors"

	"google.golang.org/grpc/grpclog"

	"github.com/vectorman1/analysis/analysis-api/common"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jackc/pgx"

	"github.com/vectorman1/analysis/analysis-api/generated/instrument_service"
)

type InstrumentsServiceContract interface {
	// repo methods
	Get(ctx context.Context, uuid string) (*instrument_service.Instrument, error)
	GetPaged(ctx context.Context, req *instrument_service.PagedRequest) (*[]*instrument_service.Instrument, uint, error)
	Overview(ctx context.Context, req *instrument_service.InstrumentRequest) (*instrument_service.InstrumentOverview, error)

	// service methods
	UpdateAll(ctx context.Context) (*instrument_service.UpdateAllResponse, error)
	recalculateRelevantInstruments(
		input []*instrument_service.InstrumentStatus,
		ctx context.Context) (*instrument_service.UpdateAllResponse, error)
	symbolDataToEntity(in *[]*instrument_service.Instrument) ([]*model.Instrument, error)
	filterUnusableSymbols(symbols *[]*instrument_service.Instruments) *[]*instrument_service.Instrument
}

type InstrumentsService struct {
	symbolRepository         *repo.InstrumentRepository
	symbolOverviewRepository *repo.OverviewRepository
	alphaVantageService      *third_party.AlphaVantageService
	externalSymbolService    *third_party.ExternalSymbolService
}

func NewSymbolService(
	symbolsRepository *repo.InstrumentRepository,
	symbolOverviewRepository *repo.OverviewRepository,
	alphaVantageService *third_party.AlphaVantageService,
	externalSymbolService *third_party.ExternalSymbolService) *InstrumentsService {
	return &InstrumentsService{
		symbolRepository:         symbolsRepository,
		symbolOverviewRepository: symbolOverviewRepository,
		alphaVantageService:      alphaVantageService,
		externalSymbolService:    externalSymbolService,
	}
}

func (s *InstrumentsService) Get(ctx context.Context, uuid string) (*instrument_service.Instrument, error) {
	sym, err := s.symbolRepository.GetByUuid(ctx, uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid symbol uuid")
	}

	return sym.ToProto(), nil
}

func (s *InstrumentsService) GetPaged(
	ctx context.Context,
	req *instrument_service.PagedRequest) (*[]*instrument_service.Instrument, uint, error) {
	if req.Filter == nil {
		return nil, 0, status.Errorf(codes.InvalidArgument, "provide filter")
	}
	if req.Filter.Order == "" {
		return nil, 0, status.Error(codes.InvalidArgument, "provide order argument")
	}

	var res []*instrument_service.Instrument
	syms, totalItemsCount, err := s.symbolRepository.GetPaged(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	for _, sym := range *syms {
		res = append(res, sym.ToProto())
	}

	return &res, totalItemsCount, nil
}

func (s *InstrumentsService) Overview(
	ctx context.Context,
	req *instrument_service.InstrumentRequest) (*instrument_service.InstrumentOverview, error) {
	userInfo := ctx.Value("user_info")
	if userInfo == nil {
		return nil, status.Error(codes.Unauthenticated, "provide user token")
	}

	symbol, err := s.symbolRepository.GetByUuid(ctx, req.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, validationErrors.NoSymbolFound)
	}
	psym := symbol.ToProto()

	overview, err := s.symbolOverviewRepository.GetByInstrumentUuid(ctx, req.Uuid)
	if err != nil {
		newOverview, err := s.getAndInsertInstrumentOverview(ctx, psym)
		if err != nil {
			return nil, status.Error(codes.NotFound, validationErrors.NoOverviewFoundForSymbol)
		}

		overview = newOverview
	} else if overview.ShouldUpdate() {
		err := s.symbolOverviewRepository.Delete(ctx, psym.Uuid)
		if err != nil {
			return nil, err
		}

		newOverview, err := s.getAndInsertInstrumentOverview(ctx, psym)
		if err != nil {
			return nil, status.Error(codes.NotFound, validationErrors.NoOverviewFoundForSymbol)
		}

		overview = newOverview
	}

	return overview.ToProto(), nil
}

func (s *InstrumentsService) UpdateAll(ctx context.Context) (*instrument_service.UpdateAllResponse, error) {
	oldSymbols, _, err := s.GetPaged(ctx, &instrument_service.PagedRequest{
		Filter: &instrument_service.PagedFilter{
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
	response, err := s.recalculateRelevantInstruments(result, ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *InstrumentsService) recalculateRelevantInstruments(
	input []*instrument_service.InstrumentStatus,
	ctx context.Context) (*instrument_service.UpdateAllResponse, error) {

	var createSymbols []*instrument_service.Instrument
	var updateSymbols []*instrument_service.Instrument
	var deleteSymbols []*instrument_service.Instrument

	itemsIgnored := int64(0)
	for _, res := range input {
		switch res.Type {
		case instrument_service.InstrumentStatus_CREATE:
			createSymbols = append(createSymbols, res.Symbol)
		case instrument_service.InstrumentStatus_UPDATE:
			updateSymbols = append(updateSymbols, res.Symbol)
		case instrument_service.InstrumentStatus_DELETE:
			deleteSymbols = append(deleteSymbols, res.Symbol)
		case instrument_service.InstrumentStatus_IGNORE:
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

	return &instrument_service.UpdateAllResponse{
		ItemsCreated: int64(len(createSymbols)),
		ItemsUpdated: int64(len(updateSymbols)),
		ItemsDeleted: int64(len(deleteSymbols)),
		ItemsIgnored: itemsIgnored,
		TotalItems:   int64(len(input)),
	}, nil
}

func (s *InstrumentsService) symbolDataToEntity(in *[]*instrument_service.Instrument) ([]*model.Instrument, error) {
	var result []*model.Instrument

	for _, sym := range *in {
		result = append(result, model.Instrument{}.FromProtoObject(sym))
	}

	return result, nil
}

func (s *InstrumentsService) generateRecalculationResult(
	newSymbols []*instrument_service.Instrument,
	oldSymbols []*instrument_service.Instrument) []*instrument_service.InstrumentStatus {
	var result []*instrument_service.InstrumentStatus
	unique := make(map[string]*instrument_service.InstrumentStatus)
	output := make(chan *instrument_service.InstrumentStatus)

	go func() {
		for _, newSym := range newSymbols {
			if ok, oldSym := common.ContainsSymbol(newSym.Uuid, oldSymbols); !ok {
				output <- &instrument_service.InstrumentStatus{
					Type:   instrument_service.InstrumentStatus_CREATE,
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
					output <- &instrument_service.InstrumentStatus{
						Type:   instrument_service.InstrumentStatus_UPDATE,
						Symbol: newSym,
					}
					continue
				} else {
					output <- &instrument_service.InstrumentStatus{
						Type:   instrument_service.InstrumentStatus_IGNORE,
						Symbol: newSym,
					}
				}
			}
		}

		for _, sym := range oldSymbols {
			if ok, _ := common.ContainsSymbol(sym.Uuid, newSymbols); !ok {
				output <- &instrument_service.InstrumentStatus{
					Type:   instrument_service.InstrumentStatus_DELETE,
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

func (s *InstrumentsService) filterUnusableSymbols(symbols *[]*instrument_service.Instrument) *[]*instrument_service.Instrument {
	var res []*instrument_service.Instrument

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

func (s *InstrumentsService) getAndInsertInstrumentOverview(ctx context.Context, sym *instrument_service.Instrument) (*model.Overview, error) {
	extOverview, err := s.alphaVantageService.GetInstrumentOverview(sym.Identifier)
	if err != nil {
		return nil, err
	}

	_, err = s.symbolOverviewRepository.Insert(ctx, extOverview.ToEntity(sym.Uuid))
	if err != nil {
		return nil, err
	}

	newOverview, err := s.symbolOverviewRepository.GetByInstrumentUuid(ctx, sym.Uuid)
	if err != nil {
		return nil, err
	}

	return newOverview, nil
}
