BEGIN;

CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address VARCHAR(42) NOT NULL,
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    decimals INTEGER NOT NULL CHECK (decimals >= 0 AND decimals <= 255),
    total_supply NUMERIC(78, 0) NOT NULL,
    chain VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS token_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_id UUID NOT NULL,
    token_address VARCHAR(42) NOT NULL,
    price DOUBLE PRECISION NOT NULL DEFAULT 0,
    price_change_24h DOUBLE PRECISION NOT NULL DEFAULT 0,
    volume_24h DOUBLE PRECISION NOT NULL DEFAULT 0,
    market_cap DOUBLE PRECISION NOT NULL DEFAULT 0,
    holders BIGINT NOT NULL DEFAULT 0,
    transactions_24h BIGINT NOT NULL DEFAULT 0,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

-- Constraints
ALTER TABLE token_metrics ADD CONSTRAINT fk_token_metrics_token FOREIGN KEY (token_id) REFERENCES tokens(id);

-- Unique
ALTER TABLE tokens ADD CONSTRAINT uk_tokens_address_chain UNIQUE (address, chain);

COMMIT;