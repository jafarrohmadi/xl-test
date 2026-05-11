# XLSmart API Contract Index (3 Service Boundaries)

Primary source of truth for API contracts and sequence diagrams aligned to:
- **Order Service** (SF Backend)
- **Payment Service** (SF Payment)
- **Partner Integration** (Partner)

`Web/App` is an actor, not a service boundary.

---

## Service Contracts

### 1) Order Service (SF Backend)
| Resource | Link |
|---|---|
| OpenAPI Specification | [sf-backend/sf-backend-api.yaml](sf-backend/sf-backend-api.yaml) |
| Sequence Diagrams (per endpoint/case) | [sf-backend/sequence.md](sf-backend/sequence.md) |

**Endpoints**
- `POST /orders` - Create new order (orchestrates full flow)
- `POST /orders/fulfillment/callback` - Receive fulfillment callback from partner
- `POST /internal/notifications` - Trigger notification (internal only)

**Responsibility**: Order orchestration, source of truth for order lifecycle

### 2) Payment Service (SF Payment)
| Resource | Link |
|---|---|
| OpenAPI Specification | [sf-payment/sf-payment-api.yaml](sf-payment/sf-payment-api.yaml) |
| Sequence Diagrams (per endpoint/case) | [sf-payment/sequence.md](sf-payment/sequence.md) |

**Endpoints**
- `POST /payments` - Process payment request
- `POST /payments/webhook` - Handle payment gateway webhook

**Responsibility**: Payment processing, webhook handling, payment status management

### 3) Partner Integration (Partner)
| Resource | Link |
|---|---|
| OpenAPI Specification | [partner/partner-api.yaml](partner/partner-api.yaml) |
| Sequence Diagrams (per endpoint/case) | [partner/sequence.md](partner/sequence.md) |

**Endpoints**
- `POST /partners/orders` - Submit order to partner system
- `POST /partners/fulfillment` - Request voucher fulfillment

**Responsibility**: External partner communication, order submission, fulfillment coordination

**Note**: Partner sends fulfillment callback to Order Service endpoint `/orders/fulfillment/callback`

---

## 🔄 API Refactoring Summary

### What Changed (v2.0.0 → v3.0.0)

**Naming Improvements**:
- ✅ Removed service prefixes (`/sf-backend`, `/sf-payment`, `/partner`)
- ✅ Removed unnecessary verbs (`/submit`, `/request`, `/trigger`)
- ✅ Applied kebab-case consistently
- ✅ Domain-first structure (`/orders`, `/payments`, `/partners`)

**Old → New Endpoint Mapping**:

| Old Endpoint | New Endpoint | Service |
|--------------|--------------|---------|
| `POST /sf-backend/orders/submit` | `POST /orders` | Order Service |
| `POST /sf-backend/fulfillment/callback` | `POST /orders/fulfillment/callback` | Order Service |
| `POST /sf-backend/internal/notifications/trigger` | `POST /internal/notifications` | Order Service |
| `POST /sf-payment/payments/request` | `POST /payments` | Payment Service |
| `POST /sf-payment/payments/webhook` | `POST /payments/webhook` | Payment Service |
| `POST /partner/orders/submit` | `POST /partners/orders` | Partner Integration |
| `POST /partner/orders/fulfillment` | `POST /partners/fulfillment` | Partner Integration |

**Removed Endpoints**:
- `GET /orders/{orderId}` - Not in current scope
- `GET /payments/{referenceId}` - Not in current scope

**Service Boundaries Clarified**:
- **Order Service**: Order orchestration and lifecycle management
- **Payment Service**: Payment processing and webhook handling
- **Partner Integration**: External partner communication

---

## Architecture & Shared Components

| Resource | Description |
|---|---|
| [c4-context.md](c4-context.md) | C4 context of SF Backend, SF Payment, Partner, and callback/notification routes |
| [shared-components.yaml](shared-components.yaml) | Shared bearer auth, idempotency/signature headers, and common response envelopes |

---

## Security Baseline

All contracts must apply:
- Bearer auth for client-facing/private endpoints.
- `Idempotency-Key` for submit and callback endpoints.
- `X-Signature` + `X-Signature-Timestamp` for callback/webhook authenticity and replay protection.
- `X-Request-Id` for audit tracing.
- No secrets in examples, and sanitized error payloads.

---

## Usage Notes

Load each service file directly into Swagger or Postman:
- `sf-backend/sf-backend-api.yaml` (Order Service)
- `sf-payment/sf-payment-api.yaml` (Payment Service)
- `partner/partner-api.yaml` (Partner Integration)

Bundle with refs from `api-contract/` root when needed:

```bash
npx @redocly/cli bundle sf-backend/sf-backend-api.yaml -o dist/order-service-bundled.yaml
npx @redocly/cli bundle sf-payment/sf-payment-api.yaml -o dist/payment-service-bundled.yaml
npx @redocly/cli bundle partner/partner-api.yaml -o dist/partner-integration-bundled.yaml
```

---

## Migration Guide

### For API Consumers

**Breaking Changes**:
1. All endpoint paths have changed - update your API client configurations
2. Some response codes have been clarified (e.g., `ORDER_SUBMITTED` → `ORDER_CREATED`)

**Non-Breaking Changes**:
1. Request/response schemas remain the same
2. Authentication mechanisms unchanged
3. Idempotency and signature validation unchanged

**Recommended Migration Steps**:
1. Update base URLs and endpoint paths in your API clients
2. Test in sandbox environment first
3. Update error code handling if needed
4. Deploy to production after validation

### For API Providers

**Implementation Changes Required**:
1. Update route handlers to match new endpoint paths
2. Update response code constants
3. Update API documentation and examples
4. Maintain backward compatibility during transition period (if needed)
