# Order Service - Sequence Diagrams (Per Endpoint)

## Endpoint: `POST /orders`

```mermaid
sequenceDiagram
    participant WA as Web/Apps
    participant OS as Order Service
    participant DB as Database
    participant PI as Partner Integration
    participant PS as Payment Service

    WA->>WA: User opens sfshop.id and chooses product
    WA->>OS: POST /orders<br/>(orderId, partnerId, goods[sku, name, desc, qty], totalPrice)
    OS->>DB: Validate params & record into DB
    DB-->>OS: Order saved

    OS->>PI: POST /partners/orders<br/>(referenceId, orderId, goods, totalPrice)
    PI->>PI: Receive & process order submission
    PI-->>OS: 200 PARTNER_ORDER_ACCEPTED<br/>(partnerOrderId, status)

    OS->>OS: Update order with partnerOrderId

    OS->>PS: POST /payments<br/>(referenceId, price)
    PS->>PS: Process payment
    PS-->>OS: 200 PAYMENT_SUCCESS<br/>(referenceId, paymentStatus, paidAt)

    OS->>OS: Update order payment status

    OS->>PI: POST /partners/fulfillment<br/>(referenceId, partnerOrderId)
    PI->>PI: Process fulfillment voucher (async)
    PI-->>OS: 200 FULFILLMENT_IN_PROGRESS

    OS-->>WA: 201 ORDER_CREATED<br/>(orderId, referenceId, partnerOrderId, paymentStatus)

    Note over PI,OS: Partner will asynchronously call back<br/>POST /orders/fulfillment/callback<br/>when voucher fulfillment is ready
```

---

## Endpoint: `POST /orders/fulfillment/callback`

> Called by **Partner Integration** (not by Order Service itself) when fulfillment is complete.

```mermaid
sequenceDiagram
    participant PI as Partner Integration
    participant OS as Order Service
    participant DB as Database
    participant NW as Notification Worker
    participant WA as Web/Apps

    PI->>OS: POST /orders/fulfillment/callback<br/>(referenceId, partnerOrderId, status, voucher)
    Note over OS: Validate HMAC signature<br/>Check replay window (< 5 min)<br/>Check idempotency key

    OS->>DB: Save voucher + update order status
    DB-->>OS: Updated

    OS-->>PI: 200 FULFILLMENT_CALLBACK_ACCEPTED

    OS->>NW: POST /internal/notifications<br/>(event_type, reference_id, user_id, channels, data)
    NW->>NW: Enqueue notification job
    NW-->>OS: 202 NOTIFICATION_QUEUED
    NW->>WA: Deliver notification (email/sms/push)
    Note over WA: User receives fulfillment notification
```

---

## Endpoint: `POST /internal/notifications`

> Internal only — triggered by Order Service after fulfillment callback is processed.

```mermaid
sequenceDiagram
    participant OS as Order Service
    participant NW as Notification Worker
    participant NP as Notification Provider
    participant WA as Web/Apps

    OS->>NW: POST /internal/notifications<br/>(event_type, reference_id, user_id, channels, data)
    Note over NW: Internal endpoint only<br/>Validate request_id and payload
    NW->>NW: Enqueue notification job
    NW-->>OS: 202 NOTIFICATION_QUEUED
    NW->>NP: Send via email/sms/push provider
    NP->>WA: Deliver notification
    Note over WA: User receives notification
```
