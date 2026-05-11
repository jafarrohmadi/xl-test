# XLSmart API - C4 Context (3 Service Boundaries)

This document reflects the refactored service boundaries with clean API design:
- **Order Service** (SF Backend)
- **Payment Service** (SF Payment)
- **Partner Integration** (Partner)

Web/App remains an external actor.

## System Context

```mermaid
C4Context
    title System Context - XLSmart 3-Service Model (v3.0.0)

    Person(webApp, "Web/App Client", "User opens sfshop.id, chooses product, submits order, receives notifications")

    System_Boundary(coreBoundary, "XLSmart Core") {
        Container(orderService, "Order Service", "Order Orchestration", "POST /orders: Validates params, records to DB, orchestrates partner submission, payment, and fulfillment. POST /orders/fulfillment/callback: Receives partner callback, saves voucher, triggers notification. POST /internal/notifications: Triggers async notification processing")
        Container(paymentService, "Payment Service", "Payment Processing", "POST /payments: Processes payment requests. POST /payments/webhook: Handles payment gateway webhook with signature validation")
        Container(notificationWorker, "Notification Worker", "Async Processor", "POST /internal/notifications: Consumes notification event and sends to notification provider")
    }

    System_Boundary(partnerBoundary, "Partner Boundary") {
        Container(partnerIntegration, "Partner Integration", "External Partner Communication", "POST /partners/orders: Receives order, returns partnerOrderId. POST /partners/fulfillment: Processes voucher fulfillment. Sends callback to /orders/fulfillment/callback with voucher")
        Container(notificationProvider, "Notification Provider", "Email/SMS/Push Provider", "External channel for message delivery to end user")
    }

    System_Ext(paymentGateway, "Payment Gateway", "External payment gateway sending signed webhook callbacks")
    ContainerDb(primaryDb, "Primary Database", "PostgreSQL", "Order and fulfillment records")
    ContainerDb(paymentDb, "Payment Database", "PostgreSQL", "Payment transactions and webhook audit")

    Rel(webApp, orderService, "1. Create order (orderId, partnerId, goods, totalPrice)", "HTTPS/JSON + Bearer")
    Rel(orderService, primaryDb, "2. Validate params & record into DB", "SQL")
    Rel(orderService, partnerIntegration, "3. Submit order to partner (referenceId, orderId, goods, totalPrice)", "HTTPS/JSON")
    Rel(partnerIntegration, orderService, "4. Return partnerOrderId", "HTTPS/JSON")
    Rel(orderService, paymentService, "5. Process payment (referenceId, price)", "HTTPS/JSON")
    Rel(paymentService, paymentDb, "Persist payment transactions", "SQL")
    Rel(paymentService, orderService, "6. Response payment status", "HTTPS/JSON")
    Rel(orderService, partnerIntegration, "7. Request fulfillment (referenceId, partnerOrderId)", "HTTPS/JSON")
    Rel(partnerIntegration, orderService, "8. Fulfillment callback (referenceId, partnerOrderId, status, voucher)", "HTTPS/JSON + HMAC")
    Rel(orderService, notificationWorker, "9. Trigger notification (event_type, reference_id, user_id, channels)", "Internal Event")
    Rel(notificationWorker, notificationProvider, "Send email/sms/push", "HTTPS/API")
    Rel(notificationProvider, webApp, "10. Sent notif to user", "Push/Email/SMS")
    Rel(paymentGateway, paymentService, "Payment webhook callback", "HTTPS/JSON + HMAC")
```

## Responsibility Split

### Order Service (SF Backend)
- **POST /orders**: Orchestrates order creation flow
  - Validates params & records into DB
  - Submits order to partner
  - Processes payment
  - Requests fulfillment
- **POST /orders/fulfillment/callback**: Receives partner callback
  - Validates signature, replay window, idempotency
  - Saves voucher + updates status
  - Triggers notification
- **POST /internal/notifications**: Internal notification trigger endpoint

### Payment Service (SF Payment)
- **POST /payments**: Processes payment request from Order Service
  - Validates payload, idempotency key
  - Executes payment processing
  - Returns payment status
- **POST /payments/webhook**: Handles payment gateway webhook
  - Verifies HMAC signature, replay window
  - Updates payment status
  - Notifies Order Service of payment result

### Partner Integration (Partner)
- **POST /partners/orders**: Receives & processes order submission
  - Validates payload, idempotency key
  - Creates partnerOrderId
  - Returns partnerOrderId
- **POST /partners/fulfillment**: Processes fulfillment request
  - Validates reference and eligibility
  - Processes voucher fulfillment async
  - Returns FULFILLMENT_IN_PROGRESS
  - Sends callback to Order Service `/orders/fulfillment/callback` when ready

## API Design Principles Applied

### ✅ Clean Naming
- Kebab-case only (`/orders`, `/payments`, `/partners`)
- No verbs in endpoint names (removed `/submit`, `/request`, `/trigger`)
- Domain-first structure (no service prefixes like `/sf-backend`)

### ✅ Clear Service Boundaries
- **Order Service** = orchestration and source of truth
- **Payment Service** = payment processing only
- **Partner Integration** = external communication only
- **Internal Service** = async operations only

### ✅ Explicit Callbacks
- `/webhook` for external payment gateway callbacks
- `/callback` for external partner callbacks
- `/internal` for internal-only endpoints

### ✅ RESTful Design
- `POST /orders` - Create resource
- `GET /orders/{orderId}` - Read resource
- `GET /payments/{referenceId}` - Read resource
- Proper HTTP methods and status codes
