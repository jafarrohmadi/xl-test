-- SF Backend (Order Service) Database Schema
-- Database: sf_backend_db

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id VARCHAR(100) UNIQUE NOT NULL,
    reference_id VARCHAR(100) UNIQUE NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    partner_id VARCHAR(100) NOT NULL,
    partner_order_id VARCHAR(100),
    total_price DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    payment_status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_reference_id ON orders(reference_id);
CREATE INDEX idx_orders_status ON orders(status);

-- Order items table
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Fulfillments table (Store result from partner callbacks)
CREATE TABLE IF NOT EXISTS fulfillments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    reference_id VARCHAR(100) NOT NULL,
    partner_order_id VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    voucher_code VARCHAR(255),
    voucher_serial_number VARCHAR(255),
    failure_reason TEXT,
    fulfilled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fulfillments_reference_id ON fulfillments(reference_id);

-- Idempotency records table
CREATE TABLE IF NOT EXISTS idempotency_records (
    idempotency_key VARCHAR(100) PRIMARY KEY,
    endpoint VARCHAR(255) NOT NULL,
    status_code INTEGER NOT NULL,
    response_body JSONB,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
