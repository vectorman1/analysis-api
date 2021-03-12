package alpha_vantage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vectorman1/analysis/analysis-api/model/service"

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

func (s *AlphaVantageService) GetSymbolOverview(symbolName string) (*service.SymbolOverview, error) {
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

	var result *service.SymbolOverview
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
