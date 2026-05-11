package request_test

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/xlsmart-api/partner/model/request"
)

var validate = validator.New()

// ────────────────────────────────────────────────────────────────────────────
// PartnerGoods
// ────────────────────────────────────────────────────────────────────────────

func TestPartnerGoods_JSONSerialization(t *testing.T) {
	goods := request.PartnerGoods{
		SKU: "SKU-001",
		Qty: 3,
	}

	data, err := json.Marshal(goods)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if decoded["sku"] != "SKU-001" {
		t.Errorf("PartnerGoods sku = %v, want SKU-001", decoded["sku"])
	}
	qty, _ := decoded["qty"].(float64)
	if int(qty) != 3 {
		t.Errorf("PartnerGoods qty = %v, want 3", decoded["qty"])
	}
}

func TestPartnerGoods_Validation_Valid(t *testing.T) {
	goods := request.PartnerGoods{
		SKU: "SKU-001",
		Qty: 1,
	}

	if err := validate.Struct(goods); err != nil {
		t.Errorf("expected valid PartnerGoods to pass validation, got: %v", err)
	}
}

func TestPartnerGoods_Validation_MissingSKU(t *testing.T) {
	goods := request.PartnerGoods{
		SKU: "",
		Qty: 1,
	}

	err := validate.Struct(goods)
	if err == nil {
		t.Error("expected validation error for missing SKU, got nil")
	}
}

func TestPartnerGoods_Validation_ZeroQty(t *testing.T) {
	goods := request.PartnerGoods{
		SKU: "SKU-001",
		Qty: 0,
	}

	err := validate.Struct(goods)
	if err == nil {
		t.Error("expected validation error for qty=0 (min=1), got nil")
	}
}

func TestPartnerGoods_Validation_NegativeQty(t *testing.T) {
	goods := request.PartnerGoods{
		SKU: "SKU-001",
		Qty: -1,
	}

	err := validate.Struct(goods)
	if err == nil {
		t.Error("expected validation error for negative qty, got nil")
	}
}

// ────────────────────────────────────────────────────────────────────────────
// PartnerSubmitRequest
// ────────────────────────────────────────────────────────────────────────────

func TestPartnerSubmitRequest_JSONSerialization(t *testing.T) {
	req := request.PartnerSubmitRequest{
		ReferenceID: "REF-001",
		OrderID:     "ORD-001",
		Goods: []request.PartnerGoods{
			{SKU: "SKU-001", Qty: 2},
			{SKU: "SKU-002", Qty: 1},
		},
		TotalPrice: 150000.50,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if decoded["referenceId"] != "REF-001" {
		t.Errorf("PartnerSubmitRequest referenceId = %v, want REF-001", decoded["referenceId"])
	}
	if decoded["orderId"] != "ORD-001" {
		t.Errorf("PartnerSubmitRequest orderId = %v, want ORD-001", decoded["orderId"])
	}

	goods, ok := decoded["goods"].([]interface{})
	if !ok || len(goods) != 2 {
		t.Errorf("PartnerSubmitRequest goods = %v, want 2 items", decoded["goods"])
	}

	totalPrice, _ := decoded["totalPrice"].(float64)
	if totalPrice != 150000.50 {
		t.Errorf("PartnerSubmitRequest totalPrice = %v, want 150000.50", decoded["totalPrice"])
	}
}

func TestPartnerSubmitRequest_Validation_Valid(t *testing.T) {
	req := request.PartnerSubmitRequest{
		ReferenceID: "REF-001",
		OrderID:     "ORD-001",
		Goods:       []request.PartnerGoods{{SKU: "SKU-001", Qty: 1}},
		TotalPrice:  50000,
	}

	if err := validate.Struct(req); err != nil {
		t.Errorf("expected valid PartnerSubmitRequest, got: %v", err)
	}
}

func TestPartnerSubmitRequest_Validation_MissingReferenceID(t *testing.T) {
	req := request.PartnerSubmitRequest{
		ReferenceID: "",
		OrderID:     "ORD-001",
		Goods:       []request.PartnerGoods{{SKU: "SKU-001", Qty: 1}},
		TotalPrice:  50000,
	}

	if err := validate.Struct(req); err == nil {
		t.Error("expected validation error for missing referenceId")
	}
}

func TestPartnerSubmitRequest_Validation_MissingOrderID(t *testing.T) {
	req := request.PartnerSubmitRequest{
		ReferenceID: "REF-001",
		OrderID:     "",
		Goods:       []request.PartnerGoods{{SKU: "SKU-001", Qty: 1}},
		TotalPrice:  50000,
	}

	if err := validate.Struct(req); err == nil {
		t.Error("expected validation error for missing orderId")
	}
}

func TestPartnerSubmitRequest_Validation_EmptyGoods(t *testing.T) {
	req := request.PartnerSubmitRequest{
		ReferenceID: "REF-001",
		OrderID:     "ORD-001",
		Goods:       []request.PartnerGoods{},
		TotalPrice:  50000,
	}

	if err := validate.Struct(req); err == nil {
		t.Error("expected validation error for empty goods slice (min=1)")
	}
}

func TestPartnerSubmitRequest_Validation_NilGoods(t *testing.T) {
	req := request.PartnerSubmitRequest{
		ReferenceID: "REF-001",
		OrderID:     "ORD-001",
		Goods:       nil,
		TotalPrice:  50000,
	}

	if err := validate.Struct(req); err == nil {
		t.Error("expected validation error for nil goods")
	}
}

func TestPartnerSubmitRequest_Validation_ZeroTotalPrice(t *testing.T) {
	// totalPrice: validate:"required,min=0" – zero passes because min=0 is satisfied
	req := request.PartnerSubmitRequest{
		ReferenceID: "REF-001",
		OrderID:     "ORD-001",
		Goods:       []request.PartnerGoods{{SKU: "SKU-001", Qty: 1}},
		TotalPrice:  0,
	}

	// totalPrice=0 with "required" should fail for numeric types in go-playground/validator
	// (required treats zero as invalid for numerics)
	if err := validate.Struct(req); err == nil {
		t.Error("expected validation error for totalPrice=0 with 'required' tag")
	}
}

func TestPartnerSubmitRequest_Validation_DiveIntoGoods(t *testing.T) {
	// Goods has a valid parent but invalid child (missing SKU)
	req := request.PartnerSubmitRequest{
		ReferenceID: "REF-001",
		OrderID:     "ORD-001",
		Goods:       []request.PartnerGoods{{SKU: "", Qty: 1}},
		TotalPrice:  50000,
	}

	if err := validate.Struct(req); err == nil {
		t.Error("expected validation error for goods item with empty SKU (dive validation)")
	}
}

func TestPartnerSubmitRequest_JSONRoundTrip(t *testing.T) {
	original := request.PartnerSubmitRequest{
		ReferenceID: "REF-RT-001",
		OrderID:     "ORD-RT-001",
		Goods: []request.PartnerGoods{
			{SKU: "SKU-RT-001", Qty: 3},
		},
		TotalPrice: 75000,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var restored request.PartnerSubmitRequest
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if restored.ReferenceID != original.ReferenceID {
		t.Errorf("ReferenceID round-trip: got %v, want %v", restored.ReferenceID, original.ReferenceID)
	}
	if restored.OrderID != original.OrderID {
		t.Errorf("OrderID round-trip: got %v, want %v", restored.OrderID, original.OrderID)
	}
	if len(restored.Goods) != len(original.Goods) {
		t.Errorf("Goods length round-trip: got %v, want %v", len(restored.Goods), len(original.Goods))
	}
	if restored.Goods[0].SKU != original.Goods[0].SKU {
		t.Errorf("Goods[0].SKU round-trip: got %v, want %v", restored.Goods[0].SKU, original.Goods[0].SKU)
	}
	if restored.TotalPrice != original.TotalPrice {
		t.Errorf("TotalPrice round-trip: got %v, want %v", restored.TotalPrice, original.TotalPrice)
	}
}

// ────────────────────────────────────────────────────────────────────────────
// FulfillmentRequest
// ────────────────────────────────────────────────────────────────────────────

func TestFulfillmentRequest_JSONSerialization(t *testing.T) {
	req := request.FulfillmentRequest{
		ReferenceID:    "REF-001",
		PartnerOrderID: "P-abc123",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if decoded["referenceId"] != "REF-001" {
		t.Errorf("FulfillmentRequest referenceId = %v, want REF-001", decoded["referenceId"])
	}
	if decoded["partnerOrderId"] != "P-abc123" {
		t.Errorf("FulfillmentRequest partnerOrderId = %v, want P-abc123", decoded["partnerOrderId"])
	}
}

func TestFulfillmentRequest_Validation_Valid(t *testing.T) {
	req := request.FulfillmentRequest{
		ReferenceID:    "REF-001",
		PartnerOrderID: "P-abc123",
	}

	if err := validate.Struct(req); err != nil {
		t.Errorf("expected valid FulfillmentRequest, got: %v", err)
	}
}

func TestFulfillmentRequest_Validation_MissingReferenceID(t *testing.T) {
	req := request.FulfillmentRequest{
		ReferenceID:    "",
		PartnerOrderID: "P-abc123",
	}

	if err := validate.Struct(req); err == nil {
		t.Error("expected validation error for missing referenceId")
	}
}

func TestFulfillmentRequest_Validation_MissingPartnerOrderID(t *testing.T) {
	req := request.FulfillmentRequest{
		ReferenceID:    "REF-001",
		PartnerOrderID: "",
	}

	if err := validate.Struct(req); err == nil {
		t.Error("expected validation error for missing partnerOrderId")
	}
}

func TestFulfillmentRequest_Validation_BothMissing(t *testing.T) {
	req := request.FulfillmentRequest{
		ReferenceID:    "",
		PartnerOrderID: "",
	}

	err := validate.Struct(req)
	if err == nil {
		t.Error("expected validation error when both fields missing")
	}

	valErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	if len(valErrs) < 2 {
		t.Errorf("expected at least 2 validation errors, got %d", len(valErrs))
	}
}

func TestFulfillmentRequest_JSONRoundTrip(t *testing.T) {
	original := request.FulfillmentRequest{
		ReferenceID:    "REF-RT-001",
		PartnerOrderID: "P-rt9876",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var restored request.FulfillmentRequest
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if restored.ReferenceID != original.ReferenceID {
		t.Errorf("ReferenceID round-trip: got %v, want %v", restored.ReferenceID, original.ReferenceID)
	}
	if restored.PartnerOrderID != original.PartnerOrderID {
		t.Errorf("PartnerOrderID round-trip: got %v, want %v", restored.PartnerOrderID, original.PartnerOrderID)
	}
}
