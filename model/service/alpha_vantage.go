package service

import (
	"time"

	"github.com/jackc/pgx/pgtype"
	dbmodel "github.com/vectorman1/analysis/analysis-api/model/db"
)

type SymbolOverview struct {
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

func (s *SymbolOverview) ToEntity(uuid string) *dbmodel.SymbolOverview {
	var u pgtype.UUID
	u.Set(uuid)

	return &dbmodel.SymbolOverview{
		SymbolUuid:                 u,
		Description:                s.Description,
		Country:                    s.Country,
		Sector:                     s.Sector,
		Industry:                   s.Industry,
		Address:                    s.Address,
		FullTimeEmployees:          s.FullTimeEmployees,
		FiscalYearEnd:              s.FiscalYearEnd,
		LatestQuarter:              s.LatestQuarter,
		MarketCapitalization:       s.MarketCapitalization,
		EBITDA:                     s.EBITDA,
		PERatio:                    s.PERatio,
		PEGRatio:                   s.PEGRatio,
		BookValue:                  s.BookValue,
		DividendPerShare:           s.DividendPerShare,
		DividendYield:              s.DividendYield,
		EPS:                        s.EPS,
		RevenuePerShareTTM:         s.RevenuePerShareTTM,
		ProfitMargin:               s.ProfitMargin,
		OperatingMarginTTM:         s.OperatingMarginTTM,
		ReturnOnAssetsTTM:          s.ReturnOnAssetsTTM,
		ReturnOnEquityTTM:          s.ReturnOnEquityTTM,
		RevenueTTM:                 s.RevenueTTM,
		GrossProfitTTM:             s.GrossProfitTTM,
		DilutedEPSTTM:              s.DilutedEPSTTM,
		QuarterlyEarningsGrowthYOY: s.QuarterlyEarningsGrowthYOY,
		QuarterlyRevenueGrowthYOY:  s.QuarterlyRevenueGrowthYOY,
		AnalystTargetPrice:         s.AnalystTargetPrice,
		TrailingPE:                 s.TrailingPE,
		ForwardPE:                  s.ForwardPE,
		PriceToSalesRatioTTM:       s.PriceToSalesRatioTTM,
		PriceToBookRatio:           s.PriceToBookRatio,
		EVToRevenue:                s.EVToRevenue,
		EVToEBITDA:                 s.EVToEBITDA,
		Beta:                       s.Beta,
		WeekHigh52:                 s.WeekLow52,
		WeekLow52:                  s.WeekLow52,
		SharesOutstanding:          s.SharesOutstanding,
		SharesFloat:                s.SharesFloat,
		SharesShort:                s.SharesShort,
		SharesShortPriorMonth:      s.SharesShortPriorMonth,
		ShortRatio:                 s.ShortRatio,
		ShortPercentOutstanding:    s.ShortPercentOutstanding,
		ShortPercentFloat:          s.ShortPercentFloat,
		PercentInsiders:            s.PercentInsiders,
		PercentInstitutions:        s.PercentInstitutions,
		ForwardAnnualDividendRate:  s.ForwardAnnualDividendRate,
		ForwardAnnualDividendYield: s.ForwardAnnualDividendYield,
		PayoutRatio:                s.PayoutRatio,
		DividendDate:               s.DividendDate,
		ExDividendDate:             s.ExDividendDate,
		LastSplitFactor:            s.LastSplitFactor,
		LastSplitDate:              s.LastSplitDate,
		UpdatedAt: pgtype.Timestamptz{
			Time:   time.Now(),
			Status: pgtype.Present,
		},
	}
}
