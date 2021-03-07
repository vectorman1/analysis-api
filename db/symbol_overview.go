package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/jackc/pgx"
	dbmodel "github.com/vectorman1/analysis/analysis-api/model/db"
)

type symbolOverviewRepository interface {
	Insert(ctx context.Context, overview *dbmodel.SymbolOverview) (bool, error)
	GetBySymbolUuid(ctx context.Context, uuid string) (*dbmodel.SymbolOverview, error)
}

type SymbolOverviewRepository struct {
	symbolOverviewRepository
	db *pgx.ConnPool
}

func NewSymbolOverviewRepository(db *pgx.ConnPool) *SymbolOverviewRepository {
	return &SymbolOverviewRepository{
		db: db,
	}
}

func (r *SymbolOverviewRepository) Insert(ctx context.Context, overview *dbmodel.SymbolOverview) (bool, error) {
	queryBuilder := squirrel.
		Insert("analysis.symbol_overviews").
		Columns("symbol_uuid, description, country, sector, industry, address, full_time_employees, fiscal_year_end, latest_quarter, market_capitalization, ebitda, pe_ratio, peg_ratio, book_value, dividend_per_share, dividend_yield, eps, revenue_per_share_ttm, profit_margin, operating_margin_ttm, return_on_assets_ttm, return_on_equity_ttm, revenue_ttm, gross_profit_ttm, diluted_eps_ttm, quarterly_earnings_growth_yoy, quarterly_revenue_growth_yoy, analyst_target_price, trailing_pe, forward_pe, price_to_sales_ratio_ttm, price_to_book_ratio, ev_to_revenue, ev_to_ebitda, beta, week_high_52, week_low_52, shares_outstanding, shares_float, shares_short, shares_short_prior_month, short_ratio, short_percent_outstanding, short_percent_float, percent_insiders, percent_institutions, forward_annual_dividend_date, forward_annual_dividend_yield, payout_ratio, dividend_date, ex_dividend_date, last_split_factor, last_split_date, updated_at").
		Values(
			&overview.SymbolUuid,
			&overview.Description,
			&overview.Country,
			&overview.Sector,
			&overview.Industry,
			&overview.Address,
			&overview.FullTimeEmployees,
			&overview.FiscalYearEnd,
			&overview.LatestQuarter,
			&overview.MarketCapitalization,
			&overview.EBITDA,
			&overview.PERatio,
			&overview.PEGRatio,
			&overview.BookValue,
			&overview.DividendPerShare,
			&overview.DividendYield,
			&overview.EPS,
			&overview.RevenuePerShareTTM,
			&overview.ProfitMargin,
			&overview.OperatingMarginTTM,
			&overview.ReturnOnAssetsTTM,
			&overview.ReturnOnEquityTTM,
			&overview.RevenueTTM,
			&overview.GrossProfitTTM,
			&overview.DilutedEPSTTM,
			&overview.QuarterlyEarningsGrowthYOY,
			&overview.QuarterlyRevenueGrowthYOY,
			&overview.AnalystTargetPrice,
			&overview.TrailingPE,
			&overview.ForwardPE,
			&overview.PriceToSalesRatioTTM,
			&overview.PriceToBookRatio,
			&overview.EVToRevenue,
			&overview.EVToEBITDA,
			&overview.Beta,
			&overview.WeekHigh52,
			&overview.WeekLow52,
			&overview.SharesOutstanding,
			&overview.SharesFloat,
			&overview.SharesShort,
			&overview.SharesShortPriorMonth,
			&overview.ShortRatio,
			&overview.ShortPercentOutstanding,
			&overview.ShortPercentFloat,
			&overview.PercentInsiders,
			&overview.PercentInstitutions,
			&overview.ForwardAnnualDividendRate,
			&overview.ForwardAnnualDividendYield,
			&overview.PayoutRatio,
			&overview.DividendDate,
			&overview.ExDividendDate,
			&overview.LastSplitFactor,
			&overview.LastSplitDate,
			&overview.UpdatedAt,
		).
		PlaceholderFormat(squirrel.Dollar)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return false, err
	}

	_, err = r.db.ExecEx(ctx, query, &pgx.QueryExOptions{}, args...)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *SymbolOverviewRepository) GetBySymbolUuid(ctx context.Context, uuid string) (*dbmodel.SymbolOverview, error) {
	queryBuilder := squirrel.
		Select("*").
		From("analysis.symbol_overviews").
		Where(fmt.Sprintf("symbol_uuid = '%s'", uuid)).
		Limit(1)
	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var result dbmodel.SymbolOverview
	row := r.db.QueryRowEx(ctx, query, &pgx.QueryExOptions{})
	if err = row.Scan(
		&result.SymbolUuid,
		&result.Description,
		&result.Country,
		&result.Sector,
		&result.Industry,
		&result.Address,
		&result.FullTimeEmployees,
		&result.FiscalYearEnd,
		&result.LatestQuarter,
		&result.MarketCapitalization,
		&result.EBITDA,
		&result.PERatio,
		&result.PEGRatio,
		&result.BookValue,
		&result.DividendPerShare,
		&result.DividendYield,
		&result.EPS,
		&result.RevenuePerShareTTM,
		&result.ProfitMargin,
		&result.OperatingMarginTTM,
		&result.ReturnOnAssetsTTM,
		&result.ReturnOnEquityTTM,
		&result.RevenueTTM,
		&result.GrossProfitTTM,
		&result.DilutedEPSTTM,
		&result.QuarterlyEarningsGrowthYOY,
		&result.QuarterlyRevenueGrowthYOY,
		&result.AnalystTargetPrice,
		&result.TrailingPE,
		&result.ForwardPE,
		&result.PriceToSalesRatioTTM,
		&result.PriceToBookRatio,
		&result.EVToRevenue,
		&result.EVToEBITDA,
		&result.Beta,
		&result.WeekHigh52,
		&result.WeekLow52,
		&result.SharesOutstanding,
		&result.SharesFloat,
		&result.SharesShort,
		&result.SharesShortPriorMonth,
		&result.ShortRatio,
		&result.ShortPercentOutstanding,
		&result.ShortPercentFloat,
		&result.PercentInsiders,
		&result.PercentInstitutions,
		&result.ForwardAnnualDividendRate,
		&result.ForwardAnnualDividendYield,
		&result.PayoutRatio,
		&result.DividendDate,
		&result.ExDividendDate,
		&result.LastSplitFactor,
		&result.LastSplitDate,
		&result.UpdatedAt); err != nil {
		return nil, err
	}

	return &result, nil
}
