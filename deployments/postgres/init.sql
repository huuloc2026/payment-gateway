CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL,
    amount NUMERIC(18,2) NOT NULL CHECK (amount >= 0),
    status VARCHAR(16) NOT NULL,
    signature VARCHAR(256),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
