package aggregate

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestOrder_JSON(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	
	order := Order{
		ID:            id,
		OrderID:       "ORD-001",
		ReferenceID:   "REF-001",
		UserID:        "user-1",
		PartnerID:     "partner-1",
		TotalPrice:    150.50,
		Status:        "PENDING",
		PaymentStatus: "UNPAID",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("failed to marshal order: %v", err)
	}

	var unmarshaled Order
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal order: %v", err)
	}

	if unmarshaled.ID != order.ID {
		t.Errorf("expected ID %v, got %v", order.ID, unmarshaled.ID)
	}
	if unmarshaled.OrderID != order.OrderID {
		t.Errorf("expected OrderID %s, got %s", order.OrderID, unmarshaled.OrderID)
	}
	if unmarshaled.TotalPrice != order.TotalPrice {
		t.Errorf("expected TotalPrice %f, got %f", order.TotalPrice, unmarshaled.TotalPrice)
	}
}

func TestOrderItem_JSON(t *testing.T) {
	id := uuid.New()
	orderID := uuid.New()
	
	item := OrderItem{
		ID:          id,
		OrderID:     orderID,
		SKU:         "SKU-123",
		Name:        "Test Item",
		Quantity:    2,
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("failed to marshal order item: %v", err)
	}

	var unmarshaled OrderItem
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal order item: %v", err)
	}

	if unmarshaled.SKU != item.SKU {
		t.Errorf("expected SKU %s, got %s", item.SKU, unmarshaled.SKU)
	}
	if unmarshaled.Quantity != item.Quantity {
		t.Errorf("expected Quantity %d, got %d", item.Quantity, unmarshaled.Quantity)
	}
}

func TestFulfillment_JSON(t *testing.T) {
	id := uuid.New()
	orderID := uuid.New()
	
	fulfillment := Fulfillment{
		ID:             id,
		OrderID:        orderID,
		ReferenceID:    "REF-001",
		PartnerOrderID: "P-ORD-001",
		Status:         "SUCCESS",
		VoucherData:    `{"code": "VOUCHER123"}`,
	}

	data, err := json.Marshal(fulfillment)
	if err != nil {
		t.Fatalf("failed to marshal fulfillment: %v", err)
	}

	var unmarshaled Fulfillment
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal fulfillment: %v", err)
	}

	if unmarshaled.Status != fulfillment.Status {
		t.Errorf("expected Status %s, got %s", fulfillment.Status, unmarshaled.Status)
	}
}
