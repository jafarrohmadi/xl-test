# SF Payment - Sequence Diagrams (Per Endpoint)

## Endpoint: `POST /payments`

```mermaid
sequenceDiagram
    participant SF as SF Backend
    participant SFP as SF Payment
    participant DB as PaymentDB

    SF->>SFP: POST /payments<br/>(referenceId, price)
    Note over SFP: Process Payment<br/>Validate payload, idempotency key
    SFP->>DB: Insert payment status PENDING
    DB-->>SFP: persisted
    SFP->>SFP: Execute payment processing
    SFP->>DB: Update status SUCCESS + paidAt
    SFP-->>SF: 200 PAYMENT_SUCCESS<br/>(referenceId, paymentStatus, paidAt)
```

## Endpoint: `POST /payments/webhook`

```mermaid
sequenceDiagram
    participant GW as Payment Gateway
    participant SFP as SF Payment
    participant DB as PaymentDB
    participant SF as SF Backend

    GW->>SFP: POST /payments/webhook<br/>(referenceId, transactionId, status, amount)
    Note over SFP: Verify HMAC signature, replay window, idempotency key
    SFP->>DB: Load payment by referenceId
    DB-->>SFP: payment record
    SFP->>DB: Update payment status from webhook
    DB-->>SFP: updated
    SFP->>SF: NotifyPaymentResult(referenceId, status)
    SF-->>SFP: accepted
    SFP-->>GW: 200 WEBHOOK_ACCEPTED
```
