# Partner (Partner Integration Service)

Partner integration microservice - handles communication with external partner systems.

## 🎯 Responsibility

- Order submission to external partners
- Fulfillment request management
- Partner API integration
- Callback sending to SF Backend

## 📡 API Endpoints

### Internal API (Called by SF Backend)
- `POST /partners/orders` - Submit order to partner system
- `POST /partners/fulfillment` - Request voucher fulfillment

### Outbound Callbacks
- Sends callback to SF Backend: `POST http://sf-backend:8081/orders/fulfillment/callback`

## 🗄️ Database

**Database Name**: `partner_db`

**Tables**:
- `partner_orders` - Orders submitted to partners
- `partner_order_items` - Order line items
- `partner_fulfillments` - Fulfillment requests and status
- `partner_configurations` - Partner system configurations
- `partner_api_logs` - API call audit log
- `callback_queue` - Callback queue for async processing
- `idempotency_records` - Idempotency tracking

## 🔗 Service Dependencies

### Outbound (Calls to other services)
- **External Partner APIs**: Various partner systems
- **SF Backend**: `POST /orders/fulfillment/callback` - Send fulfillment callback

### Inbound (Receives calls from)
- **SF Backend**: `POST /partners/orders` - Order submission
- **SF Backend**: `POST /partners/fulfillment` - Fulfillment request

## 🚀 Running Locally

### Prerequisites
- Go 1.21+
- PostgreSQL 15+

### Environment Variables
```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5434/partner_db?sslmode=disable

# Server
PORT=8083
ENV=development

# Partner API Config
PARTNER_A_API_URL=https://api.partner-a.example.com
PARTNER_A_API_KEY=your-partner-a-api-key
PARTNER_A_SECRET=your-partner-a-secret

# Callback Config
SF_BACKEND_CALLBACK_URL=http://localhost:8081/orders/fulfillment/callback

# Security
HMAC_SECRET=your-hmac-secret
```

### Run
```bash
# Install dependencies
go mod download

# Run database migrations
psql -U postgres -d partner_db -f database/partner.sql

# Run service
go run cmd/main.go
```

### Docker
```bash
# Build
docker build -t partner:latest .

# Run
docker run -p 8083:8083 \
  -e DATABASE_URL=postgres://postgres:postgres@host.docker.internal:5434/partner_db \
  partner:latest
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
curl http://localhost:8083/health
```

## 🔒 Security

- **Idempotency**: Prevent duplicate order submission
- **HMAC Signature**: Sign callbacks to SF Backend
- **API Key Management**: Secure partner credentials

## 📝 Example Request

### Submit Order to Partner
```bash
curl -X POST http://localhost:8083/partners/orders \
  -H "Content-Type: application/json" \
  -H "X-Request-Id: req-123" \
  -H "Idempotency-Key: idem-456" \
  -d '{
    "referenceId": "REF-001",
    "orderId": "ORD-001",
    "goods": [
      {
        "sku": "SKU-001",
        "qty": 2
      }
    ],
    "totalPrice": 150000
  }'
```

### Request Fulfillment
```bash
curl -X POST http://localhost:8083/partners/fulfillment \
  -H "Content-Type: application/json" \
  -H "X-Request-Id: req-789" \
  -H "Idempotency-Key: idem-789" \
  -d '{
    "referenceId": "REF-001",
    "partnerOrderId": "P-778899"
  }'
```