-- Partner (Partner Integration Service) Database Schema
-- Database: partner_db

-- Partner orders table
CREATE TABLE IF NOT EXISTS partner_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference_id VARCHAR(100) UNIQUE NOT NULL,
    order_id VARCHAR(100) NOT NULL,
    partner_order_id VARCHAR(100) UNIQUE NOT NULL,
    partner_id VARCHAR(100) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL CHECK (total_price >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    submitted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_partner_orders_reference_id ON partner_orders(reference_id);

-- Partner order items table (Simplified)
CREATE TABLE IF NOT EXISTS partner_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_order_id UUID NOT NULL REFERENCES partner_orders(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Partner fulfillments table
CREATE TABLE IF NOT EXISTS partner_fulfillments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_order_id UUID NOT NULL REFERENCES partner_orders(id) ON DELETE CASCADE,
    reference_id VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    voucher_code VARCHAR(255),
    voucher_serial_number VARCHAR(255),
    failure_reason TEXT,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fulfilled_at TIMESTAMP,
    callback_sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
