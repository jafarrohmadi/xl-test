# Order Service - Sequence Diagrams (Per Endpoint)

## Endpoint: `POST /orders`

```mermaid
sequenceDiagram
    participant WebApps as Web/Apps
    participant Backend as sf-backend
    participant Database as Database
    participant Partner as partner
    participant Payment as sf-payment

    WebApps->>WebApps: User opens sfshop.id and chooses product
    WebApps->>Backend: POST /orders<br/>(orderId, partnerId, goods[sku, name, desc, qty], totalPrice)
    Backend->>Database: SELECT from idempotency_records
    Database-->>Backend: idempotency result
    
    alt Idempotency Hit
        Backend-->>WebApps: 200 OK (Return cached response)
    else Idempotency Miss
        Backend->>Database: INSERT into orders (Status: PENDING)
        Database-->>Backend: order saved
        Backend->>Database: INSERT into order_items
        Database-->>Backend: items saved

        Backend->>Partner: POST /partners/orders<br/>(referenceId, orderId, goods, totalPrice)
        Partner->>Partner: Receive & process order submission
        
        alt Partner Success
            Partner-->>Backend: 200 PARTNER_ORDER_ACCEPTED
            Backend->>Database: UPDATE orders (Update partnerOrderId)
            Database-->>Backend: updated
        else Partner Failure
            Partner-->>Backend: 400/500 Error
            Backend->>Database: UPDATE orders (Status: FAILED)
            Backend-->>WebApps: Error Response
        end
    end

    Backend->>Payment: POST /payments<br/>(referenceId, price)
    Payment->>Payment: Process payment
    Payment-->>Backend: 200 PAYMENT_SUCCESS<br/>(referenceId, paymentStatus, paidAt)

    Backend->>Database: UPDATE orders (Update payment status)
    Database-->>Backend: updated

    Backend->>Partner: POST /partners/fulfillment<br/>(referenceId, partnerOrderId)
    Partner->>Partner: Process fulfillment voucher (async)
    Partner-->>Backend: 200 FULFILLMENT_IN_PROGRESS

    Backend->>Database: UPDATE orders (Status: SUBMITTED)
    Database-->>Backend: updated
    Backend-->>WebApps: 201 ORDER_CREATED<br/>(orderId, referenceId, partnerOrderId, paymentStatus)

    Note over Partner,Backend: Partner will asynchronously call back<br/>POST /orders/fulfillment/callback<br/>when voucher fulfillment is ready
```

---

## Endpoint: `POST /orders/fulfillment/callback`

> Called by **Partner Integration** (not by Order Service itself) when fulfillment is complete.

```mermaid
sequenceDiagram
    participant Partner as partner
    participant Backend as sf-backend
    participant Database as Database
    participant NotificationWorker as Notification Worker
    participant WebApps as Web/Apps

    Partner->>Backend: POST /orders/fulfillment/callback<br/>(referenceId, partnerOrderId, status, voucher)
    Note over Backend: Validate HMAC signature<br/>Check replay window (< 5 min)<br/>Check idempotency key

    Backend->>Database: SELECT from orders (by ReferenceID)
    Database-->>Backend: order record
    Backend->>Database: INSERT into fulfillments (Save voucher data)
    Database-->>Backend: fulfillment saved
    Backend->>Database: UPDATE orders (Status: COMPLETED/FAILED)
    Database-->>Backend: updated

    Backend-->>Partner: 200 FULFILLMENT_CALLBACK_ACCEPTED

    Backend->>NotificationWorker: POST /internal/notifications<br/>(event_type, reference_id, user_id, channels, data)
    NotificationWorker->>NotificationWorker: Enqueue notification job
    NotificationWorker-->>Backend: 202 NOTIFICATION_QUEUED
    NotificationWorker->>WebApps: Deliver notification (email/sms/push)
    Note over WebApps: User receives fulfillment notification
```

---

## Endpoint: `POST /internal/notifications`

> Internal only — triggered by Order Service after fulfillment callback is processed.

```mermaid
sequenceDiagram
    participant Backend as sf-backend
    participant NotificationWorker as Notification Worker
    participant NotificationProvider as Notification Provider
    participant WebApps as Web/Apps

    Backend->>NotificationWorker: POST /internal/notifications<br/>(event_type, reference_id, user_id, channels, data)
    Note over NotificationWorker: Internal endpoint only<br/>Validate request_id and payload
    NotificationWorker->>NotificationWorker: Enqueue notification job
    NotificationWorker-->>Backend: 202 NOTIFICATION_QUEUED
    NotificationWorker->>NotificationProvider: Send via email/sms/push provider
    NotificationProvider->>WebApps: Deliver notification
    Note over WebApps: User receives notification
```
