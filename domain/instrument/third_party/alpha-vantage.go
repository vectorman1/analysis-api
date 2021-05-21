package third_party

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/model"

	validationErrors "github.com/vectorman1/analysis/analysis-api/common/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vectorman1/analysis/analysis-api/common"
)

const (
	SymbolOverviewEndpoint = "https://www.alphavantage.co/query?function=OVERVIEW&symbol=%s&apikey=%s"
)

type alphaVantageService interface {
}

type AlphaVantageService struct {
	alphaVantageService
	httpClient *http.Client
	config     *common.Config
}

func NewAlphaVantageService(config *common.Config) *AlphaVantageService {
	client := &http.Client{Timeout: 5 * time.Second}

	return &AlphaVantageService{
		config:     config,
		httpClient: client,
	}
}

func (s *AlphaVantageService) GetInstrumentOverview(symbolName string) (*model.InstrumentOverviewResponse, error) {
	url := fmt.Sprintf(SymbolOverviewEndpoint, symbolName, s.config.AlphaVantageApiKey)
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	res, err := s.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result *model.InstrumentOverviewResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Description == "" {
		return nil, status.Error(codes.NotFound, validationErrors.NoOverviewFoundForSymbol)
	}

	return result, nil
}
