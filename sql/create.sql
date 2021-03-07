CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS analysis;

DO
$do$
BEGIN
        IF NOT EXISTS (
                SELECT FROM pg_catalog.pg_user
                WHERE usename = 'harb') THEN
            CREATE USER harb WITH ENCRYPTED PASSWORD 'HueHue123';
END IF;
END
$do$;

GRANT ALL PRIVILEGES ON DATABASE analysis TO harb;
GRANT ALL PRIVILEGES ON SCHEMA analysis TO harb;

CREATE TABLE IF NOT EXISTS analysis.currencies
(
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    long_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS analysis.symbols
(
    id SERIAL PRIMARY KEY,
    uuid uuid NOT NULL,
    currency_id BIGINT NOT NULL,
    isin TEXT NOT NULL,
    identifier TEXT NOT NULL,
    name TEXT NOT NULL,
    minimum_order_quantity REAL NOT NULL,
    market_name TEXT NOT NULL,
    market_hours_gmt TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
    CONSTRAINT fk_symbols_currency FOREIGN KEY (currency_id) REFERENCES analysis.currencies (id) ON UPDATE NO ACTION ON DELETE NO ACTION
    );

CREATE TABLE IF NOT EXISTS analysis.histories
(
    id SERIAL PRIMARY KEY,
    symbol_uuid uuid NOT NULL,
    values TEXT NOT NULL,
    for_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_histories_symbol FOREIGN KEY (symbol_uuid) REFERENCES analysis.symbols (uuid) ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS analysis.symbol_overviews
(
    symbol_uuid uuid PRIMARY KEY NOT NULL,
    description TEXT NOT NULL,
    country TEXT NOT NULL,
    sector TEXT NOT NULL,
    industry TEXT NOT NULL,
    address TEXT NOT NULL,
    full_time_employees TEXT NOT NULL,
    fiscal_year_end TEXT NOT NULL,
    latest_quarter TEXT NOT NULL,
    market_capitalization TEXT NOT NULL,
    ebitda TEXT NOT NULL,
    pe_ratio TEXT NOT NULL,
    peg_ratio TEXT NOT NULL,
    book_value TEXT NOT NULL,
    dividend_per_share TEXT NOT NULL,
    dividend_yield TEXT NOT NULL,
    eps TEXT NOT NULL,
    revenue_per_share_ttm TEXT NOT NULL,
    profit_margin TEXT NOT NULL,
    operating_margin_ttm TEXT NOT NULL,
    return_on_assets_ttm TEXT NOT NULL,
    return_on_equity_ttm TEXT NOT NULL,
    revenue_ttm TEXT NOT NULL,
    gross_profit_ttm TEXT NOT NULL,
    diluted_eps_ttm TEXT NOT NULL,
    quarterly_earnings_growth_yoy TEXT NOT NULL,
    quarterly_revenue_growth_yoy TEXT NOT NULL,
    analyst_target_price TEXT NOT NULL,
    trailing_pe TEXT NOT NULL,
    forward_pe TEXT NOT NULL,
    price_to_sales_ratio_ttm TEXT NOT NULL,
    price_to_book_ratio TEXT NOT NULL,
    ev_to_revenue TEXT NOT NULL,
    ev_to_ebitda TEXT NOT NULL,
    beta TEXT NOT NULL,
    week_high_52 TEXT NOT NULL,
    week_low_52 TEXT NOT NULL,
    shares_outstanding TEXT NOT NULL,
    shares_float TEXT NOT NULL,
    shares_short TEXT NOT NULL,
    shares_short_prior_month TEXT NOT NULL,
    short_ratio TEXT NOT NULL,
    short_percent_outstanding TEXT NOT NULL,
    short_percent_float TEXT NOT NULL,
    percent_insiders TEXT NOT NULL,
    percent_institutions TEXT NOT NULL,
    forward_annual_dividend_date TEXT NOT NULL,
    forward_annual_dividend_yield TEXT NOT NULL,
    payout_ratio TEXT NOT NULL,
    dividend_date TEXT NOT NULL,
    ex_dividend_date TEXT NOT NULL,
    last_split_factor TEXT NOT NULL,
    last_split_date TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_overview_symbol FOREIGN KEY (symbol_uuid) REFERENCES analysis.symbols (uuid) ON UPDATE NO ACTION ON DELETE NO ACTION
)

CREATE TABLE IF NOT EXISTS analysis.reports
(
    id SERIAL PRIMARY KEY,
    symbol_id BIGINT NOT NULL,
    history_id BIGINT NOT NULL,
    exponential_moving_averages TEXT NOT NULL,
    simple_moving_averages TEXT NOT NULL,
    macd TEXT NOT NULL,
    rsi TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
    CONSTRAINT fk_report_symbol FOREIGN KEY (symbol_id) REFERENCES analysis.symbols (id) ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT fk_report_history FOREIGN KEY (history_id) REFERENCES analysis.histories (id) ON UPDATE NO ACTION ON DELETE NO ACTION
    );

CREATE TABLE IF NOT EXISTS analysis.strategies
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
    );

CREATE TABLE IF NOT EXISTS analysis.signals
(
    id SERIAL PRIMARY KEY,
    symbol_id BIGINT NOT NULL,
    strategy_id BIGINT NOT NULL,
    type BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
    CONSTRAINT fk_signal_symbol FOREIGN KEY (symbol_id) REFERENCES analysis.symbols (id) ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT fk_report_strategy FOREIGN KEY (strategy_id) REFERENCES analysis.strategies (id) ON UPDATE NO ACTION ON DELETE NO ACTION
    );

CREATE TABLE IF NOT EXISTS analysis.signal_reports
(
    id SERIAL PRIMARY KEY,
    signal_id BIGINT NOT NULL,
    report_id BIGINT NOT NULL,
    CONSTRAINT fk_signal_report_signal FOREIGN KEY (signal_id) REFERENCES analysis.signals (id) ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT fk_signal_report_report FOREIGN KEY (report_id) REFERENCES analysis.reports (id) ON UPDATE NO ACTION ON DELETE NO ACTION
    );

CREATE SCHEMA IF NOT EXISTS "user";

CREATE TABLE IF NOT EXISTS "user".users
(
    id SERIAL PRIMARY KEY,
    uuid uuid NOT NULL UNIQUE,
    private_role BIGINT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
    );
