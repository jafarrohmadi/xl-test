# XLSmart API - Microservices Architecture

Production-ready microservices architecture with 3 independent services and separate databases.

## 🏗️ Architecture Overview
test case
![diagram.png](api-contract/diagram.png)


```mermaid
C4Container
    title Container diagram for XLSmart Microservices

    Person(user, "Web/App User", "A user interacting with XLSmart apps")
    
    System_Boundary(xlsmart, "XLSmart System") {
        Container(sf_backend, "SF Backend", "Go", "Order orchestration and lifecycle management. Port: 8081")
        ContainerDb(sf_backend_db, "SF Backend DB", "PostgreSQL", "Stores orders and notifications")
        
        Container(sf_payment, "SF Payment", "Go", "Payment processing and gateway integration. Port: 8082")
        ContainerDb(sf_payment_db, "SF Payment DB", "PostgreSQL", "Stores payment transactions")
        
        Container(partner, "Partner", "Go", "External partner communication and fulfillment. Port: 8083")
        ContainerDb(partner_db, "Partner DB", "PostgreSQL", "Stores partner requests")
    }

    System_Ext(payment_gateway, "Payment Gateway", "External payment provider")
    System_Ext(external_partner, "External Partner", "Partner fulfillment system")

    Rel(user, sf_backend, "Submits orders", "REST API")
    
    Rel(sf_backend, sf_backend_db, "Reads/Writes")
    Rel(sf_payment, sf_payment_db, "Reads/Writes")
    Rel(partner, partner_db, "Reads/Writes")

    Rel(sf_backend, partner, "Submits partner order / Requests fulfillment", "REST API")
    Rel(sf_backend, sf_payment, "Processes payment", "REST API")
    
    Rel(payment_gateway, sf_payment, "Sends webhook", "REST API")
    Rel(external_partner, partner, "Processes fulfillment asynchronously", "REST API")
    Rel(partner, sf_backend, "Sends fulfillment callback", "REST API")
```

---

## 🚀 How to Run

### Prerequisites
- **Docker** & **Docker Compose**
- Postman (for testing API)

### Step-by-Step Guide (Docker)

1. **Clone and Setup**
   Make sure you are in the project root directory where `docker-compose.microservices.yml` is located.

2. **Run All Services**
   Execute the following command to start the databases and microservices:
   ```bash
   docker compose -f docker-compose.microservices.yml up -d --build
   ```

3. **Verify Status**
   Check if all 6 containers (3 Databases, 3 Services) are running successfully:
   ```bash
   docker ps
   ```

4. **View Logs (Optional)**
   If you want to see the real-time logs of the microservices:
   ```bash
   docker compose -f docker-compose.microservices.yml logs -f
   ```

5. **Stop Services**
   To turn off all services:
   ```bash
   docker compose -f docker-compose.microservices.yml down
   ```

---

## 🧪 How to Test (Postman)

A pre-configured Postman Collection is included in the project to make testing easier.

### Importing Postman Collection
1. Open your Postman application.
2. Click **Import** (top left).
3. Select or drag-and-drop the file: `xlsmart-api.postman_collection.json` located in the `api-contract/` folder of this project.
4. You will see a new collection named **"XLSmart Microservices API"**.

### Testing the Endpoints
### Synchronous Workflow Diagram (Current Implementation)

Berikut adalah urutan interaksi antar service yang berjalan saat ini menggunakan Synchronous HTTP Call:

```mermaid
sequenceDiagram
    autonumber
    actor User as Web/Apps
    participant Backend as SF Backend
    participant DB_Backend as SF Backend DB
    participant Payment as SF Payment
    participant DB_Payment as SF Payment DB
    participant Partner as Partner Service
    participant DB_Partner as Partner DB

    User->>Backend: POST /orders (Submit Order)
    activate Backend
    Backend->>DB_Backend: Check Idempotency
    Backend->>DB_Backend: Save Order (Status: PENDING)
    
    Backend->>Payment: POST /payments (Process)
    activate Payment
    Payment->>DB_Payment: Save Transaction
    Payment-->>Backend: 200 OK (Success)
    deactivate Payment
    
    Backend->>Partner: POST /partners/orders (Submit)
    activate Partner
    Partner->>DB_Partner: Save Partner Order
    Partner-->>Backend: 200 OK (Success)
    deactivate Partner

    Backend->>DB_Backend: Update Order (Status: SUBMITTED)
    Backend-->>User: 201 Created (Return Order Details)
    deactivate Backend

    Note over Backend, Partner: Async Fulfillment Process
    Partner->>Backend: POST /orders/fulfillment/callback
    activate Backend
    Backend->>DB_Backend: Update Order & Fulfillment Data
    Backend-->>Partner: 200 OK
    deactivate Backend
```

*Note: The Postman collection uses dynamic variables like `{{$guid}}` to automatically generate unique `X-Request-Id` and `Idempotency-Key` for every request.*

---

## 📦 Services Details

### 1. **SF Backend (Order Service)** - Port 8081
**Responsibility**: Order orchestration and lifecycle management
**Database**: `sf_backend_db` (PostgreSQL on port 5432)
**Endpoints**:
- `POST /orders` - Create order
- `POST /orders/fulfillment/callback` - Receive fulfillment callback
- `POST /internal/notifications` - Trigger notifications

### 2. **SF Payment (Payment Service)** - Port 8082
**Responsibility**: Payment processing and gateway integration
**Database**: `sf_payment_db` (PostgreSQL on port 5433)
**Endpoints**:
- `POST /payments` - Process payment
- `POST /payments/webhook` - Handle payment gateway webhook

### 3. **Partner (Integration Service)** - Port 8083
**Responsibility**: External partner communication
**Database**: `partner_db` (PostgreSQL on port 5434)
**Endpoints**:
- `POST /partners/orders` - Submit order to partner
- `POST /partners/fulfillment` - Request fulfillment

---

## 🔒 Security & Features
- **Data Isolation**: Each service uses its own database (`sf_backend_db`, `sf_payment_db`, `partner_db`).
- **Idempotency**: Implemented `Idempotency-Key` headers on POST endpoints to prevent double processing.
- **Traceability**: Implemented `X-Request-Id` to trace requests across microservices.
- **Signatures**: HMAC signature validation for webhooks and callbacks via `X-Signature`.

---

## 💡 Proposed Solutions for Better Improvement

While the current architecture successfully splits the domains into three independent microservices, the orchestration relies heavily on synchronous HTTP calls. Below are ideas to elevate the architecture to be more robust, scalable, and fault-tolerant.

### 1. Shift to Event-Driven Architecture (EDA) & Saga Pattern
- **Current Flaw:** The `SF Backend` synchronously calls `Partner` and `SF Payment` during the user's request. If the Payment service takes 10 seconds, the user waits 10 seconds. If `Partner` is down, the whole order fails immediately without a chance to recover.
- **Solution:** Implement the **Saga Pattern** using a Message Broker (e.g., Kafka or RabbitMQ). `SF Backend` immediately returns a `202 Accepted` to the user and processes the workflow asynchronously in the background. If a step fails (e.g., Payment fails), the Saga orchestrator will automatically trigger compensating transactions (e.g., Cancel Partner Order).

### 2. Robust Idempotency & Retry Mechanisms
- **Current State:** Basic Idempotency checks are implemented via HTTP headers (`Idempotency-Key`) and a database check.
- **Solution:** Combine this with a robust retry mechanism (Exponential Backoff) and Dead Letter Queues (DLQ) for asynchronous messages. This ensures that transient network failures do not result in lost orders or double payments.

### 3. Circuit Breaking & Rate Limiting
- **Solution:** Implement Circuit Breakers (e.g., using `gobreaker`) on inter-service HTTP clients. If the `Partner Service` starts failing or lagging, the Circuit Breaker trips and prevents cascading failures, returning a fast fallback response instead of hanging and consuming resources.

### 4. Improved Ideal Workflow (Saga Pattern)

Here is how the improved Event-Driven Orchestration (Saga) would look:

```mermaid
sequenceDiagram
    autonumber
    actor User as Web/Apps
    participant Backend as sf-backend (Orchestrator)
    participant Broker as Message Broker (Kafka)
    participant Partner as partner
    participant Payment as sf-payment

    User->>Backend: POST /orders (Submit)
    Backend->>Backend: Create Order (Status: PENDING)
    Backend-->>User: 202 Accepted

    Note over Backend, Payment: --- SAGA SUCCESS PATH ---
    
    Backend->>Broker: Command: CreatePartnerOrder
    Broker->>Partner: Consume Command
    Partner->>Partner: Process Order
    Partner-->>Broker: Reply: PartnerOrderCreated
    Broker-->>Backend: Consume Reply
    
    Backend->>Broker: Command: ProcessPayment
    Broker->>Payment: Consume Command
    Payment->>Payment: Execute Payment
    
    alt Payment Success
        Payment-->>Broker: Reply: PaymentSuccess
        Broker-->>Backend: Consume Reply
        Backend->>Broker: Command: RequestFulfillment
        Broker->>Partner: Consume Command
        Partner->>Partner: Process Fulfillment
        Partner-->>Broker: Reply: FulfillmentCompleted
        Broker-->>Backend: Consume Reply
        Backend->>Backend: Update Order (Status: COMPLETED)
        Backend->>Broker: Event: OrderCompleted
        Broker-->>User: Push Notification (WebSocket/FCM)
    else Payment Failed
        Payment-->>Broker: Reply: PaymentFailed
        Broker-->>Backend: Consume Reply
        Backend->>Broker: Command: CancelPartnerOrder (Compensation)
        Broker->>Partner: Consume Command
        Partner->>Partner: Rollback/Cancel Order
        Backend->>Backend: Update Order (Status: FAILED)
    end
```
