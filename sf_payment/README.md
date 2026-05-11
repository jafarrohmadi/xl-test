# SF Payment (Payment Service)

Payment processing microservice - handles payment requests and gateway webhooks.

## 🎯 Responsibility

- Payment request processing
- Payment gateway integration
- Webhook handling from payment gateways
- Payment status management

## 📡 API Endpoints

### Internal API
- `POST /payments` - Process payment request (called by SF Backend)

### Webhook API
- `POST /payments/webhook` - Handle payment gateway webhook

## 🗄️ Database

**Database Name**: `sf_payment_db`

**Tables**:
- `payment_transactions` - Payment transaction records
- `payment_webhooks` - Webhook audit log
- `payment_instructions` - Payment instructions (VA, QR, etc.)
- `payment_retries` - Retry tracking for failed payments
- `idempotency_records` - Idempotency tracking

## 🔗 Service Dependencies

### Outbound (Calls to other services)
- **Payment Gateway**: External payment provider API
- **SF Backend** (optional): Notify payment status updates

### Inbound (Receives calls from)
- **SF Backend**: `POST /payments` - Payment request
- **Payment Gateway**: `POST /payments/webhook` - Payment status webhook

## 🚀 Running Locally

### Prerequisites
- Go 1.21+
- PostgreSQL 15+

### Environment Variables
```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5433/sf_payment_db?sslmode=disable

# Server
PORT=8082
ENV=development

# Payment Gateway
PAYMENT_GATEWAY_URL=https://api.payment-gateway.example.com
PAYMENT_GATEWAY_API_KEY=your-api-key
PAYMENT_GATEWAY_SECRET=your-secret

# Security
HMAC_SECRET=your-hmac-secret
```

### Run
```bash
# Install dependencies
go mod download

# Run database migrations
psql -U postgres -d sf_payment_db -f database/sf_payment.sql

# Run service
go run cmd/main.go
```

### Docker
```bash
# Build
docker build -t sf-payment:latest .

# Run
docker run -p 8082:8082 \
  -e DATABASE_URL=postgres://postgres:postgres@host.docker.internal:5433/sf_payment_db \
  sf-payment:latest
```

## 🧪 Testing

```bash
# Unit tests
go test ./...

# Integration tests
go test ./... -tags=integration

# Test coverage
go test -cover ./...
```

## 📊 Health Check

```bash
curl http://localhost:8082/health
```

## 🔒 Security

- **HMAC Signature**: Webhook signature validation
- **Idempotency**: Prevent duplicate payment processing
- **Replay Protection**: Timestamp validation for webhooks

## 📝 Example Request

### Process Payment
```bash
curl -X POST http://localhost:8082/payments \
  -H "Content-Type: application/json" \
  -H "X-Request-Id: req-123" \
  -H "Idempotency-Key: idem-456" \
  -d '{
    "referenceId": "REF-001",
    "price": 150000
  }'
```

### Payment Webhook (from Gateway)
```bash
curl -X POST http://localhost:8082/payments/webhook \
  -H "Content-Type: application/json" \
  -H "X-Request-Id: req-789" \
  -H "X-Signature: hmac-signature" \
  -H "X-Signature-Timestamp: 2026-05-11T10:00:00Z" \
  -d '{
    "referenceId": "REF-001",
    "transactionId": "TXN-123",
    "status": "SUCCESS",
    "amount": 150000,
    "paidAt": "2026-05-11T10:00:00Z"
  }'
``