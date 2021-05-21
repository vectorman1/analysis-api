package model

import (
	"strconv"
	"time"
)

type InstrumentOverviewResponse struct {
	Symbol                     string `json:"Symbol"`
	AssetType                  string `json:"AssetType"`
	Name                       string `json:"Name"`
	Description                string `json:"Description"`
	Exchange                   string `json:"Exchange"`
	Currency                   string `json:"Currency"`
	Country                    string `json:"Country"`
	Sector                     string `json:"Sector"`
	Industry                   string `json:"Industry"`
	Address                    string `json:"Address"`
	FullTimeEmployees          string `json:"FullTimeEmployees"`
	FiscalYearEnd              string `json:"FiscalYearEnd"`
	LatestQuarter              string `json:"LatestQuarter"`
	MarketCapitalization       string `json:"MarketCapitalization"`
	EBITDA                     string `json:"EBITDA"`
	PERatio                    string `json:"PERatio"`
	PEGRatio                   string `json:"PEGRatio"`
	BookValue                  string `json:"BookValue"`
	DividendPerShare           string `json:"DividendPerShare"`
	DividendYield              string `json:"DividendYield"`
	EPS                        string `json:"EPS"`
	RevenuePerShareTTM         string `json:"RevenuePerShareTTM"`
	ProfitMargin               string `json:"ProfitMargin"`
	OperatingMarginTTM         string `json:"OperatingMarginTTM"`
	ReturnOnAssetsTTM          string `json:"ReturnOnAssetsTTM"`
	ReturnOnEquityTTM          string `json:"ReturnOnEquityTTM"`
	RevenueTTM                 string `json:"RevenueTTM"`
	GrossProfitTTM             string `json:"GrossProfitTTM"`
	DilutedEPSTTM              string `json:"DilutedEPSTTM"`
	QuarterlyEarningsGrowthYOY string `json:"QuarterlyEarningsGrowthYOY"`
	QuarterlyRevenueGrowthYOY  string `json:"QuarterlyRevenueGrowthYOY"`
	AnalystTargetPrice         string `json:"AnalystTargetPrice"`
	TrailingPE                 string `json:"TrailingPE"`
	ForwardPE                  string `json:"ForwardPE"`
	PriceToSalesRatioTTM       string `json:"PriceToSalesRatioTTM"`
	PriceToBookRatio           string `json:"PriceToBookRatio"`
	EVToRevenue                string `json:"EVToRevenue"`
	EVToEBITDA                 string `json:"EVToEBITDA"`
	Beta                       string `json:"Beta"`
	WeekHigh52                 string `json:"52WeekHigh"`
	WeekLow52                  string `json:"52WeekLow"`
	DayMovingAverage50         string `json:"50DayMovingAverage"`
	DayMovingAverage200        string `json:"200DayMovingAverage"`
	SharesOutstanding          string `json:"SharesOutstanding"`
	SharesFloat                string `json:"SharesFloat"`
	SharesShort                string `json:"SharesShort"`
	SharesShortPriorMonth      string `json:"SharesShortPriorMonth"`
	ShortRatio                 string `json:"ShortRatio"`
	ShortPercentOutstanding    string `json:"ShortPercentOutstanding"`
	ShortPercentFloat          string `json:"ShortPercentFloat"`
	PercentInsiders            string `json:"PercentInsiders"`
	PercentInstitutions        string `json:"PercentInstitutions"`
	ForwardAnnualDividendRate  string `json:"ForwardAnnualDividendRate"`
	ForwardAnnualDividendYield string `json:"ForwardAnnualDividendYield"`
	PayoutRatio                string `json:"PayoutRatio"`
	DividendDate               string `json:"DividendDate"`
	ExDividendDate             string `json:"ExDividendDate"`
	LastSplitFactor            string `json:"LastSplitFactor"`
	LastSplitDate              string `json:"LastSplitDate"`
}

func (s *InstrumentOverviewResponse) ToEntity(uuid string) *InstrumentOverview {
	fullTimeEmployees, _ := strconv.ParseInt(s.FullTimeEmployees, 10, 64)
	latestQuarter, _ := time.Parse("2006-01-02", s.LatestQuarter)
	marketCapitalization, _ := strconv.ParseInt(s.MarketCapitalization, 10, 64)
	ebitda, _ := strconv.ParseInt(s.EBITDA, 10, 64)
	peRatio, _ := strconv.ParseFloat(s.PERatio, 32)
	pegRatio, _ := strconv.ParseFloat(s.PEGRatio, 32)
	bookValue, _ := strconv.ParseFloat(s.PERatio, 32)
	dividendPerShare, _ := strconv.ParseFloat(s.PERatio, 32)
	dividendYield, _ := strconv.ParseFloat(s.PERatio, 32)
	revenuePerShareTTM, _ := strconv.ParseFloat(s.RevenuePerShareTTM, 32)
	profitMargin, _ := strconv.ParseFloat(s.ProfitMargin, 32)
	operatingMarginTTM, _ := strconv.ParseFloat(s.OperatingMarginTTM, 32)
	returnOnAssetsTTM, _ := strconv.ParseFloat(s.ReturnOnAssetsTTM, 32)
	eps, _ := strconv.ParseFloat(s.EPS, 32)
	returnOnEquityTTM, _ := strconv.ParseFloat(s.ReturnOnEquityTTM, 32)
	revenueTTM, _ := strconv.ParseInt(s.RevenueTTM, 10, 64)
	grossProfitTTM, _ := strconv.ParseInt(s.GrossProfitTTM, 10, 64)
	dilutedEPSTTM, _ := strconv.ParseFloat(s.DilutedEPSTTM, 32)
	quarterlyEarningsGrowthYOY, _ := strconv.ParseFloat(s.QuarterlyEarningsGrowthYOY, 32)
	quarterlyRevenueGrowthYOY, _ := strconv.ParseFloat(s.QuarterlyRevenueGrowthYOY, 32)
	analystTargetPrice, _ := strconv.ParseFloat(s.AnalystTargetPrice, 32)
	trailingPE, _ := strconv.ParseFloat(s.TrailingPE, 32)
	forwardPE, _ := strconv.ParseFloat(s.ForwardPE, 32)
	priceToSalesRatioTTM, _ := strconv.ParseFloat(s.PriceToSalesRatioTTM, 32)
	priceToBookRatio, _ := strconv.ParseFloat(s.PriceToBookRatio, 32)
	eVToRevenue, _ := strconv.ParseFloat(s.EVToRevenue, 32)
	eVToEBITDA, _ := strconv.ParseFloat(s.EVToEBITDA, 32)
	beta, _ := strconv.ParseFloat(s.Beta, 32)
	weekHigh52, _ := strconv.ParseFloat(s.WeekHigh52, 32)
	weekLow52, _ := strconv.ParseFloat(s.WeekLow52, 32)
	sharesOutstanding, _ := strconv.ParseInt(s.SharesOutstanding, 10, 32)
	sharesFloat, _ := strconv.ParseInt(s.SharesFloat, 10, 32)
	sharesShort, _ := strconv.ParseInt(s.SharesShort, 10, 32)
	sharesShortPriorMonth, _ := strconv.ParseInt(s.SharesShortPriorMonth, 10, 32)
	shortRatio, _ := strconv.ParseFloat(s.ShortRatio, 32)
	shortPercentOutstanding, _ := strconv.ParseFloat(s.ShortPercentOutstanding, 32)
	shortPercentFloat, _ := strconv.ParseFloat(s.ShortPercentFloat, 32)
	percentInsiders, _ := strconv.ParseFloat(s.PercentInsiders, 32)
	percentInstitutions, _ := strconv.ParseFloat(s.PercentInstitutions, 32)
	forwardAnnualDividendRate, _ := strconv.ParseFloat(s.ForwardAnnualDividendRate, 32)
	forwardAnnualDividendYield, _ := strconv.ParseFloat(s.ForwardAnnualDividendYield, 32)
	payoutRatio, _ := strconv.ParseFloat(s.PayoutRatio, 32)
	dividendDate, _ := time.Parse("2006-01-02", s.DividendDate)
	exDividendDate, _ := time.Parse("2006-01-02", s.ExDividendDate)
	lastSplitDate, _ := time.Parse("2006-01-02", s.LastSplitDate)

	return &InstrumentOverview{
		SymbolUuid:                 uuid,
		Description:                s.Description,
		Country:                    s.Country,
		Sector:                     s.Sector,
		Industry:                   s.Industry,
		Address:                    s.Address,
		FullTimeEmployees:          fullTimeEmployees,
		FiscalYearEnd:              s.FiscalYearEnd,
		LatestQuarter:              latestQuarter,
		MarketCapitalization:       marketCapitalization,
		EBITDA:                     ebitda,
		PERatio:                    float32(peRatio),
		PEGRatio:                   float32(pegRatio),
		BookValue:                  float32(bookValue),
		DividendPerShare:           float32(dividendPerShare),
		DividendYield:              float32(dividendYield),
		EPS:                        float32(eps),
		RevenuePerShareTTM:         float32(revenuePerShareTTM),
		ProfitMargin:               float32(profitMargin),
		OperatingMarginTTM:         float32(operatingMarginTTM),
		ReturnOnAssetsTTM:          float32(returnOnAssetsTTM),
		ReturnOnEquityTTM:          float32(returnOnEquityTTM),
		RevenueTTM:                 revenueTTM,
		GrossProfitTTM:             grossProfitTTM,
		DilutedEPSTTM:              float32(dilutedEPSTTM),
		QuarterlyEarningsGrowthYOY: float32(quarterlyEarningsGrowthYOY),
		QuarterlyRevenueGrowthYOY:  float32(quarterlyRevenueGrowthYOY),
		AnalystTargetPrice:         float32(analystTargetPrice),
		TrailingPE:                 float32(trailingPE),
		ForwardPE:                  float32(forwardPE),
		PriceToSalesRatioTTM:       float32(priceToSalesRatioTTM),
		PriceToBookRatio:           float32(priceToBookRatio),
		EVToRevenue:                float32(eVToRevenue),
		EVToEBITDA:                 float32(eVToEBITDA),
		Beta:                       float32(beta),
		WeekHigh52:                 float32(weekHigh52),
		WeekLow52:                  float32(weekLow52),
		SharesOutstanding:          sharesOutstanding,
		SharesFloat:                sharesFloat,
		SharesShort:                sharesShort,
		SharesShortPriorMonth:      sharesShortPriorMonth,
		ShortRatio:                 float32(shortRatio),
		ShortPercentOutstanding:    float32(shortPercentOutstanding),
		ShortPercentFloat:          float32(shortPercentFloat),
		PercentInsiders:            float32(percentInsiders),
		PercentInstitutions:        float32(percentInstitutions),
		ForwardAnnualDividendRate:  float32(forwardAnnualDividendRate),
		ForwardAnnualDividendYield: float32(forwardAnnualDividendYield),
		PayoutRatio:                float32(payoutRatio),
		DividendDate:               dividendDate,
		ExDividendDate:             exDividendDate,
		LastSplitFactor:            s.LastSplitFactor,
		LastSplitDate:              lastSplitDate,
		UpdatedAt:                  time.Now(),
	}
}
