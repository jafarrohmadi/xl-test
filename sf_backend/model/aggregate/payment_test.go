package aggregate

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPaymentTransaction_JSON(t *testing.T) {
	id := uuid.New()
	
	tx := PaymentTransaction{
		ID:            id,
		ReferenceID:   "REF-PAY-001",
		TransactionID: "TX-999",
		Amount:        250.75,
		Status:        "PAID",
		PaymentMethod: "OVO",
		Provider:      "XLPAY",
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("failed to marshal payment transaction: %v", err)
	}

	var unmarshaled PaymentTransaction
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal payment transaction: %v", err)
	}

	if unmarshaled.ReferenceID != tx.ReferenceID {
		t.Errorf("expected ReferenceID %s, got %s", tx.ReferenceID, unmarshaled.ReferenceID)
	}
	if unmarshaled.Amount != tx.Amount {
		t.Errorf("expected Amount %f, got %f", tx.Amount, unmarshaled.Amount)
	}
}

func TestPaymentWebhook_JSON(t *testing.T) {
	id := uuid.New()
	txID := uuid.New()
	
	webhook := PaymentWebhook{
		ID:                   id,
		PaymentTransactionID: txID,
		ReferenceID:          "REF-PAY-001",
		TransactionID:        "TX-999",
		Status:               "SUCCESS",
		Amount:               250.75,
		RawPayload:           `{"event": "payment.success"}`,
	}

	data, err := json.Marshal(webhook)
	if err != nil {
		t.Fatalf("failed to marshal payment webhook: %v", err)
	}

	var unmarshaled PaymentWebhook
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal payment webhook: %v", err)
	}

	if unmarshaled.Status != webhook.Status {
		t.Errorf("expected Status %s, got %s", webhook.Status, unmarshaled.Status)
	}
}

func TestIdempotencyRecord_JSON(t *testing.T) {
	record := IdempotencyRecord{
		IdempotencyKey: "key-123",
		Endpoint:       "/api/v1/orders",
		StatusCode:     201,
		ResponseBody:   `{"id": "ord-1"}`,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	data, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("failed to marshal idempotency record: %v", err)
	}

	var unmarshaled IdempotencyRecord
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal idempotency record: %v", err)
	}

	if unmarshaled.IdempotencyKey != record.IdempotencyKey {
		t.Errorf("expected IdempotencyKey %s, got %s", record.IdempotencyKey, unmarshaled.IdempotencyKey)
	}
}
