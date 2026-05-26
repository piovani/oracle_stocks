CREATE TABLE IF NOT EXISTS quotes (
    id              BIGSERIAL PRIMARY KEY,
    date            DATE         NOT NULL,
    ticker          VARCHAR(12)  NOT NULL,
    bdi_code        VARCHAR(2)   NOT NULL,
    market_type     VARCHAR(3)   NOT NULL,
    short_name      VARCHAR(12),
    specification   VARCHAR(10),
    currency        VARCHAR(4),
    open_price      NUMERIC(18,4),
    high_price      NUMERIC(18,4),
    low_price       NUMERIC(18,4),
    average_price   NUMERIC(18,4),
    close_price     NUMERIC(18,4),
    best_bid        NUMERIC(18,4),
    best_ask        NUMERIC(18,4),
    trades          BIGINT       NOT NULL DEFAULT 0,
    quantity        BIGINT       NOT NULL DEFAULT 0,
    volume          NUMERIC(20,2),
    strike_price    NUMERIC(18,4),
    expiry_date     DATE,
    quote_factor    BIGINT       NOT NULL DEFAULT 1,
    isin            VARCHAR(12),
    dist_number     VARCHAR(3),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_quotes_date_ticker_bdi UNIQUE (date, ticker, bdi_code)
);

CREATE INDEX IF NOT EXISTS idx_quotes_ticker_date ON quotes (ticker, date DESC);
CREATE INDEX IF NOT EXISTS idx_quotes_date        ON quotes (date);
CREATE INDEX IF NOT EXISTS idx_quotes_isin        ON quotes (isin);
