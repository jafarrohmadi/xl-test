# SF Backend (Order Service)

Order orchestration microservice - source of truth for order lifecycle.

## 🎯 Responsibility

- Order creation and management
- Order orchestration (coordinates with Payment and Partner services)
- Fulfillment callback handling
- Notification triggering

## 📡 API Endpoints

### Public API
- `POST /orders` - Create new order

### Internal Callbacks
- `POST /orders/fulfillment/callback` - Receive fulfillment callback from Partner service

### Internal Triggers
- `POST /internal/notifications` - Trigger notification (internal only)

## 🗄️ Database

**Database Name**: `sf_backend_db`

**Tables**:
- `orders` - Main orders table
- `order_items` - Order line items
- `fulfillments` - Fulfillment records
- `notifications` - Notification queue
- `idempotency_records` - Idempotency tracking
- `webhook_signatures` - Webhook signature audit

## 🔗 Service Dependencies

### Outbound (Calls to other services)
- **SF Payment Service**: `POST /payments` - Process payment
- **Partner Service**: `POST /partners/orders` - Submit order to partner
- **Partner Service**: `POST /partners/fulfillment` - Request fulfillment

### Inbound (Receives calls from)
- **Web/App**: `POST /orders` - User creates order
- **Partner Service**: `POST /orders/fulfillment/callback` - Fulfillment callback

## 🚀 Running Locally

### Prerequisites
- Go 1.21+
- PostgreSQL 15+

### Environment Variables
```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/sf_backend_db?sslmode=disable

# Service URLs
SF_PAYMENT_URL=http://localhost:8082
PARTNER_URL=http://localhost:8083

# Server
PORT=8081
ENV=development

# Security
JWT_SECRET=your-jwt-secret
HMAC_SECRET=your-hmac-secret
```

### Run
```bash
# Install dependencies
go mod download

# Run database migrations
psql -U postgres -d sf_backend_db -f database/sf_backend.sql

# Run service
go run cmd/main.go
```

### Docker
```bash
# Build
docker build -t sf-backend:latest .

# Run
docker run -p 8081:8081 \
  -e DATABASE_URL=postgres://postgres:postgres@host.docker.internal:5432/sf_backend_db \
  sf-backend:latest
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
curl http://localhost:8081/health
```

## 🔒 Security

- **Authentication**: Bearer token for public endpoints
- **HMAC Signature**: For webhook/callback validation
- **Idempotency**: All POST endpoints support idempotency keys

## 📝 Example Request

### Create Order
```bash
curl -X POST http://localhost:8081/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -H "X-Request-Id: req-123" \
  -H "Idempotency-Key: idem-456" \
  -d '{
    "orderId": "ORD-001",
    "partnerId": "PARTNER-A",
    "goods": [
      {
        "sku": "SKU-001",
        "name": "Product A",
        "qty": 2
      }
    ],
    "totalPrice": 150000
  }'
```