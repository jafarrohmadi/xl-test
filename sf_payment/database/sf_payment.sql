-- SF Payment (Payment Service) Database Schema
-- Database: sf_payment_db

-- Payment transactions table
CREATE TABLE IF NOT EXISTS payment_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference_id VARCHAR(100) UNIQUE NOT NULL,
    transaction_id VARCHAR(100) UNIQUE,
    amount DECIMAL(15,2) NOT NULL CHECK (amount >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    payment_method VARCHAR(50),
    provider VARCHAR(100),
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_transactions_reference_id ON payment_transactions(reference_id);
CREATE INDEX idx_payment_transactions_status ON payment_transactions(status);

-- Payment webhooks table
CREATE TABLE IF NOT EXISTS payment_webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_transaction_id UUID REFERENCES payment_transactions(id) ON DELETE SET NULL,
    reference_id VARCHAR(100) NOT NULL,
    transaction_id VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    raw_payload JSONB NOT NULL,
    signature TEXT NOT NULL,
    webhook_timestamp TIMESTAMP NOT NULL,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Idempotency records table
CREATE TABLE IF NOT EXISTS idempotency_records (
    idempotency_key VARCHAR(100) PRIMARY KEY,
    endpoint VARCHAR(255) NOT NULL,
    status_code INTEGER NOT NULL,
    response_body JSONB,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
