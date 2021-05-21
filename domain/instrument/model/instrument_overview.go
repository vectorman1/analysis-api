package model

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
)

type InstrumentOverview struct {
	SymbolUuid                 string
	Description                string
	Country                    string
	Sector                     string
	Industry                   string
	Address                    string
	FullTimeEmployees          int64
	FiscalYearEnd              string
	LatestQuarter              time.Time
	MarketCapitalization       int64
	EBITDA                     int64
	PERatio                    float32
	PEGRatio                   float32
	BookValue                  float32
	DividendPerShare           float32
	DividendYield              float32
	EPS                        float32
	RevenuePerShareTTM         float32
	ProfitMargin               float32
	OperatingMarginTTM         float32
	ReturnOnAssetsTTM          float32
	ReturnOnEquityTTM          float32
	RevenueTTM                 int64
	GrossProfitTTM             int64
	DilutedEPSTTM              float32
	QuarterlyEarningsGrowthYOY float32
	QuarterlyRevenueGrowthYOY  float32
	AnalystTargetPrice         float32
	TrailingPE                 float32
	ForwardPE                  float32
	PriceToSalesRatioTTM       float32
	PriceToBookRatio           float32
	EVToRevenue                float32
	EVToEBITDA                 float32
	Beta                       float32
	WeekHigh52                 float32
	WeekLow52                  float32
	SharesOutstanding          int64
	SharesFloat                int64
	SharesShort                int64
	SharesShortPriorMonth      int64
	ShortRatio                 float32
	ShortPercentOutstanding    float32
	ShortPercentFloat          float32
	PercentInsiders            float32
	PercentInstitutions        float32
	ForwardAnnualDividendRate  float32
	ForwardAnnualDividendYield float32
	PayoutRatio                float32
	DividendDate               time.Time
	ExDividendDate             time.Time
	LastSplitFactor            string
	LastSplitDate              time.Time
	UpdatedAt                  time.Time
}

func (s *InstrumentOverview) ShouldUpdate() bool {
	return time.Since(s.UpdatedAt) > 7*(24*time.Hour)
}

func (s *InstrumentOverview) ToProto() *symbol_service.SymbolOverview {
	latestQuarter := timestamppb.New(s.LatestQuarter)
	updatedAt := timestamppb.New(s.UpdatedAt)
	dividendDate := timestamppb.New(s.DividendDate)
	exDividendDate := timestamppb.New(s.ExDividendDate)
	lastSplitDate := timestamppb.New(s.DividendDate)

	return &symbol_service.SymbolOverview{
		Description:                s.Description,
		Country:                    s.Country,
		Sector:                     s.Sector,
		Industry:                   s.Industry,
		Address:                    s.Address,
		FullTimeEmployees:          s.FullTimeEmployees,
		FiscalYearEnd:              s.FiscalYearEnd,
		LatestQuarter:              latestQuarter,
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
		DividendDate:               dividendDate,
		ExDividendDate:             exDividendDate,
		LastSplitFactor:            s.LastSplitFactor,
		LastSplitDate:              lastSplitDate,
		UpdatedAt:                  updatedAt,
	}
}
