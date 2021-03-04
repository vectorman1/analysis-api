package db

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/jackc/pgx/pgtype"
	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
)

type SymbolOverview struct {
	SymbolUuid                 pgtype.UUID
	Description                string
	Country                    string
	Sector                     string
	Industry                   string
	Address                    string
	FullTimeEmployees          string
	FiscalYearEnd              string
	LatestQuarter              string
	MarketCapitalization       string
	EBITDA                     string
	PERatio                    string
	PEGRatio                   string
	BookValue                  string
	DividendPerShare           string
	DividendYield              string
	EPS                        string
	RevenuePerShareTTM         string
	ProfitMargin               string
	OperatingMarginTTM         string
	ReturnOnAssetsTTM          string
	ReturnOnEquityTTM          string
	RevenueTTM                 string
	GrossProfitTTM             string
	DilutedEPSTTM              string
	QuarterlyEarningsGrowthYOY string
	QuarterlyRevenueGrowthYOY  string
	AnalystTargetPrice         string
	TrailingPE                 string
	ForwardPE                  string
	PriceToSalesRatioTTM       string
	PriceToBookRatio           string
	EVToRevenue                string
	EVToEBITDA                 string
	Beta                       string
	WeekHigh52                 string
	WeekLow52                  string
	SharesOutstanding          string
	SharesFloat                string
	SharesShort                string
	SharesShortPriorMonth      string
	ShortRatio                 string
	ShortPercentOutstanding    string
	ShortPercentFloat          string
	PercentInsiders            string
	PercentInstitutions        string
	ForwardAnnualDividendRate  string
	ForwardAnnualDividendYield string
	PayoutRatio                string
	DividendDate               string
	ExDividendDate             string
	LastSplitFactor            string
	LastSplitDate              string
	UpdatedAt                  pgtype.Timestamptz
}

func (s *SymbolOverview) ToProtoObject() *symbol_service.SymbolOverview {
	updatedAt, _ := ptypes.TimestampProto(s.UpdatedAt.Time)

	return &symbol_service.SymbolOverview{
		Description:                s.Description,
		Country:                    s.Country,
		Sector:                     s.Sector,
		Industry:                   s.Industry,
		Address:                    s.Address,
		FullTimeEmployees:          s.FullTimeEmployees,
		FiscalYearEnd:              s.FiscalYearEnd,
		LatestQuarter:              s.LatestQuarter,
		MarketCapitalization:       s.MarketCapitalization,
		Ebitda:                     s.EBITDA,
		PeRatio:                    s.PERatio,
		PegRatio:                   s.PEGRatio,
		BookValue:                  s.BookValue,
		DividendPerShare:           s.DividendPerShare,
		DividendYield:              s.DividendYield,
		Eps:                        s.EPS,
		RevenuePerShareTtm:         s.RevenuePerShareTTM,
		ProfitMargin:               s.ProfitMargin,
		OperatingMarginTtm:         s.OperatingMarginTTM,
		ReturnOnAssetsTtm:          s.ReturnOnAssetsTTM,
		ReturnOnEquity:             s.ReturnOnEquityTTM,
		RevenueTtm:                 s.RevenueTTM,
		GrossProfitTtm:             s.GrossProfitTTM,
		DilutedEpsTtm:              s.DilutedEPSTTM,
		QuarterlyEarningsGrowthYoy: s.QuarterlyEarningsGrowthYOY,
		QuarterlyRevenueGrowthYoy:  s.QuarterlyRevenueGrowthYOY,
		AnalystTargetPrice:         s.AnalystTargetPrice,
		TrailingPe:                 s.TrailingPE,
		ForwardPe:                  s.ForwardPE,
		PriceToSalesRatioTtm:       s.PriceToSalesRatioTTM,
		PriceToBookRatio:           s.PriceToBookRatio,
		EvToRevenue:                s.EVToRevenue,
		EvToEbitda:                 s.EVToEBITDA,
		Beta:                       s.Beta,
		WeekHigh52:                 s.WeekHigh52,
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
		UpdatedAt:                  updatedAt,
	}
}
