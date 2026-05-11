# Partner - Sequence Diagrams (Per Endpoint)

## Endpoint: `POST /partners/orders`

```mermaid
sequenceDiagram
    participant Backend as sf-backend
    participant Partner as partner
    participant Database as Partner DB

    Backend->>Partner: POST /partners/orders<br/>(referenceId, orderId, goods, totalPrice)
    Partner->>Database: SELECT from idempotency_records
    Database-->>Partner: idempotency result
    
    alt Idempotency Hit
        Partner-->>Backend: 200 OK (Cached)
    else Idempotency Miss
        Partner->>Partner: Validate Stock & Price
        alt Validation Success
            Partner->>Database: INSERT into partner_orders (Status: ACCEPTED)
            Database-->>Partner: order saved
            Partner-->>Backend: 200 PARTNER_ORDER_ACCEPTED
        else Validation Failure
            Partner-->>Backend: 422 Unprocessable Entity
        end
    end
```

## Endpoint: `POST /partners/fulfillment`

```mermaid
sequenceDiagram
    participant Backend as sf-backend
    participant Partner as partner
    participant Database as Partner DB

    Backend->>Partner: POST /partners/fulfillment<br/>(referenceId, partnerOrderId)
    Partner->>Database: SELECT from partner_orders
    Database-->>Partner: partner order record
    
    alt Order Valid
        Partner->>Database: INSERT into fulfillments (Status: IN_PROGRESS)
        Database-->>Partner: fulfillment saved
        Partner->>Partner: Trigger Async Worker for Voucher
        Partner-->>Backend: 200 FULFILLMENT_IN_PROGRESS
    else Order Invalid
        Partner-->>Backend: 404 Not Found
    end
```

## Outbound Callback Contract (Partner → sf-backend)

```mermaid
sequenceDiagram
    participant Partner as partner
    participant Backend as sf-backend

    Partner->>Backend: POST /orders/fulfillment/callback<br/>(referenceId, partnerOrderId, status, voucher)
    Note over Backend: Verify HMAC signature and replay window<br/>Process callbackFulfillment
    Backend-->>Partner: 200 FULFILLMENT_CALLBACK_ACCEPTED
```
