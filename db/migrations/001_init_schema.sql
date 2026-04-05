-- Run once: psql -d finance_db -f db/migrations/001_init_schema.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT        NOT NULL,
    email      TEXT        UNIQUE NOT NULL,
    password   TEXT        NOT NULL,
    role       TEXT        NOT NULL DEFAULT 'viewer'
                           CHECK (role IN ('viewer', 'analyst', 'admin')),
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS financial_records (
    id         UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount     NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    type       TEXT           NOT NULL CHECK (type IN ('income', 'expense')),
    category   TEXT           NOT NULL,
    date       DATE           NOT NULL,
    notes      TEXT,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_records_user_id    ON financial_records(user_id);
CREATE INDEX IF NOT EXISTS idx_records_type       ON financial_records(type);
CREATE INDEX IF NOT EXISTS idx_records_category   ON financial_records(category);
CREATE INDEX IF NOT EXISTS idx_records_date       ON financial_records(date);
CREATE INDEX IF NOT EXISTS idx_records_deleted_at ON financial_records(deleted_at);
