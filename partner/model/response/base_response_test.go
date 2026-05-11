package response_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/xlsmart-api/partner/model/response"
)

// ────────────────────────────────────────────────────────────────────────────
// ResponseCode constants
// ────────────────────────────────────────────────────────────────────────────

func TestResponseCodes_SuccessHTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		code     response.ResponseCode
		wantHTTP int
		wantCode string
	}{
		{"OrderSubmitted", response.OrderSubmitted, http.StatusCreated, "ORDER_SUBMITTED"},
		{"FulfillmentCallbackAccepted", response.FulfillmentCallbackAccepted, http.StatusOK, "FULFILLMENT_CALLBACK_ACCEPTED"},
		{"NotificationQueued", response.NotificationQueued, http.StatusAccepted, "NOTIFICATION_QUEUED"},
		{"PaymentSuccess", response.PaymentSuccess, http.StatusOK, "PAYMENT_SUCCESS"},
		{"WebhookAccepted", response.WebhookAccepted, http.StatusOK, "WEBHOOK_ACCEPTED"},
		{"PartnerOrderAccepted", response.PartnerOrderAccepted, http.StatusOK, "PARTNER_ORDER_ACCEPTED"},
		{"FulfillmentInProgress", response.FulfillmentInProgress, http.StatusOK, "FULFILLMENT_IN_PROGRESS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code.HttpStatusCode != tt.wantHTTP {
				t.Errorf("%s HTTP status = %d, want %d", tt.name, tt.code.HttpStatusCode, tt.wantHTTP)
			}
			if tt.code.Code != tt.wantCode {
				t.Errorf("%s code = %q, want %q", tt.name, tt.code.Code, tt.wantCode)
			}
			if tt.code.Message == "" {
				t.Errorf("%s message should not be empty", tt.name)
			}
		})
	}
}

func TestResponseCodes_ErrorHTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		code     response.ResponseCode
		wantHTTP int
		wantCode string
	}{
		{"InvalidRequest", response.InvalidRequest, http.StatusBadRequest, "INVALID_REQUEST"},
		{"InvalidSignature", response.InvalidSignature, http.StatusBadRequest, "INVALID_SIGNATURE"},
		{"Unauthorized", response.Unauthorized, http.StatusUnauthorized, "UNAUTHORIZED"},
		{"NotFound", response.NotFound, http.StatusNotFound, "NOT_FOUND"},
		{"OrderNotFound", response.OrderNotFound, http.StatusNotFound, "ORDER_NOT_FOUND"},
		{"PaymentNotFound", response.PaymentNotFound, http.StatusNotFound, "PAYMENT_NOT_FOUND"},
		{"Conflict", response.Conflict, http.StatusConflict, "CONFLICT"},
		{"DuplicateOrder", response.DuplicateOrder, http.StatusConflict, "DUPLICATE_ORDER"},
		{"AlreadyProcessed", response.AlreadyProcessed, http.StatusConflict, "ALREADY_PROCESSED"},
		{"DuplicatePaymentRequest", response.DuplicatePaymentRequest, http.StatusConflict, "DUPLICATE_PAYMENT_REQUEST"},
		{"PartnerSubmitFailed", response.PartnerSubmitFailed, http.StatusUnprocessableEntity, "PARTNER_SUBMIT_FAILED"},
		{"PaymentFailed", response.PaymentFailed, http.StatusPaymentRequired, "PAYMENT_FAILED"},
		{"RateLimited", response.RateLimited, http.StatusTooManyRequests, "RATE_LIMITED"},
		{"InternalError", response.InternalError, http.StatusInternalServerError, "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code.HttpStatusCode != tt.wantHTTP {
				t.Errorf("%s HTTP status = %d, want %d", tt.name, tt.code.HttpStatusCode, tt.wantHTTP)
			}
			if tt.code.Code != tt.wantCode {
				t.Errorf("%s code = %q, want %q", tt.name, tt.code.Code, tt.wantCode)
			}
			if tt.code.Message == "" {
				t.Errorf("%s message should not be empty", tt.name)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────────────────────
// BuildSuccessResponse
// ────────────────────────────────────────────────────────────────────────────

func TestBuildSuccessResponse_Structure(t *testing.T) {
	resp := response.BuildSuccessResponse(response.PartnerOrderAccepted, "req-001", map[string]interface{}{
		"partnerOrderId": "P-abc123",
		"status":         "ACCEPTED",
	})

	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Code != "PARTNER_ORDER_ACCEPTED" {
		t.Errorf("code = %q, want PARTNER_ORDER_ACCEPTED", resp.Code)
	}
	if resp.RequestID != "req-001" {
		t.Errorf("request_id = %q, want req-001", resp.RequestID)
	}
	if resp.Message == "" {
		t.Error("message should not be empty")
	}
	if resp.Data == nil {
		t.Error("data should not be nil")
	}
	if resp.Details != nil {
		t.Error("details should be nil for success response")
	}
}

func TestBuildSuccessResponse_NilData(t *testing.T) {
	resp := response.BuildSuccessResponse(response.FulfillmentInProgress, "req-002", nil)

	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Code != "FULFILLMENT_IN_PROGRESS" {
		t.Errorf("code = %q, want FULFILLMENT_IN_PROGRESS", resp.Code)
	}
}

func TestBuildSuccessResponse_AllSuccessCodes(t *testing.T) {
	codes := []response.ResponseCode{
		response.OrderSubmitted,
		response.FulfillmentCallbackAccepted,
		response.NotificationQueued,
		response.PaymentSuccess,
		response.WebhookAccepted,
		response.PartnerOrderAccepted,
		response.FulfillmentInProgress,
	}

	for _, code := range codes {
		resp := response.BuildSuccessResponse(code, "req-test", nil)
		if resp == nil {
			t.Fatalf("nil response for code %v", code.Code)
		}
		if !resp.Success {
			t.Errorf("expected success=true for code %v", code.Code)
		}
		if resp.Code != code.Code {
			t.Errorf("response code mismatch: got %v, want %v", resp.Code, code.Code)
		}
	}
}

func TestBuildSuccessResponse_JSONSerialization(t *testing.T) {
	resp := response.BuildSuccessResponse(response.PartnerOrderAccepted, "req-001", map[string]interface{}{
		"partnerOrderId": "P-abc123",
	})

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if decoded["success"] != true {
		t.Errorf("JSON success = %v, want true", decoded["success"])
	}
	if decoded["code"] != "PARTNER_ORDER_ACCEPTED" {
		t.Errorf("JSON code = %v, want PARTNER_ORDER_ACCEPTED", decoded["code"])
	}
	if decoded["request_id"] != "req-001" {
		t.Errorf("JSON request_id = %v, want req-001", decoded["request_id"])
	}
	if _, exists := decoded["details"]; exists {
		t.Error("details should be omitted from JSON for success response")
	}
}

// ────────────────────────────────────────────────────────────────────────────
// BuildErrorResponse
// ────────────────────────────────────────────────────────────────────────────

func TestBuildErrorResponse_Structure(t *testing.T) {
	details := []response.ErrorDetail{{Field: "orderId", Reason: "required"}}
	resp := response.BuildErrorResponse(response.InvalidRequest, "req-001", details)

	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Success {
		t.Error("expected success=false")
	}
	if resp.Code != "INVALID_REQUEST" {
		t.Errorf("code = %q, want INVALID_REQUEST", resp.Code)
	}
	if resp.RequestID != "req-001" {
		t.Errorf("request_id = %q, want req-001", resp.RequestID)
	}
	if resp.Message == "" {
		t.Error("message should not be empty")
	}
	if resp.Data != nil {
		t.Error("data should be nil for error response")
	}
}

func TestBuildErrorResponse_NilDetails(t *testing.T) {
	resp := response.BuildErrorResponse(response.InternalError, "req-001", nil)

	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Success {
		t.Error("expected success=false")
	}
	if resp.Code != "INTERNAL_ERROR" {
		t.Errorf("code = %q, want INTERNAL_ERROR", resp.Code)
	}
}

func TestBuildErrorResponse_AllErrorCodes(t *testing.T) {
	codes := []response.ResponseCode{
		response.InvalidRequest,
		response.InvalidSignature,
		response.Unauthorized,
		response.NotFound,
		response.OrderNotFound,
		response.PaymentNotFound,
		response.Conflict,
		response.DuplicateOrder,
		response.AlreadyProcessed,
		response.DuplicatePaymentRequest,
		response.PartnerSubmitFailed,
		response.PaymentFailed,
		response.RateLimited,
		response.InternalError,
	}

	for _, code := range codes {
		resp := response.BuildErrorResponse(code, "req-test", nil)
		if resp == nil {
			t.Fatalf("nil response for code %v", code.Code)
		}
		if resp.Success {
			t.Errorf("expected success=false for code %v", code.Code)
		}
		if resp.Code != code.Code {
			t.Errorf("response code mismatch: got %v, want %v", resp.Code, code.Code)
		}
	}
}

func TestBuildErrorResponse_JSONSerialization(t *testing.T) {
	resp := response.BuildErrorResponse(response.DuplicateOrder, "req-001", nil)

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if decoded["success"] != false {
		t.Errorf("JSON success = %v, want false", decoded["success"])
	}
	if decoded["code"] != "DUPLICATE_ORDER" {
		t.Errorf("JSON code = %v, want DUPLICATE_ORDER", decoded["code"])
	}
	if decoded["request_id"] != "req-001" {
		t.Errorf("JSON request_id = %v, want req-001", decoded["request_id"])
	}
	if _, exists := decoded["data"]; exists {
		t.Error("data should be omitted from JSON for error response with nil data")
	}
}

// ────────────────────────────────────────────────────────────────────────────
// ErrorDetail
// ────────────────────────────────────────────────────────────────────────────

func TestErrorDetail_JSONSerialization(t *testing.T) {
	detail := response.ErrorDetail{
		Field:  "referenceId",
		Reason: "required",
	}

	data, err := json.Marshal(detail)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if decoded["field"] != "referenceId" {
		t.Errorf("ErrorDetail field = %v, want referenceId", decoded["field"])
	}
	if decoded["reason"] != "required" {
		t.Errorf("ErrorDetail reason = %v, want required", decoded["reason"])
	}
}

func TestErrorDetail_OmitEmptyField(t *testing.T) {
	detail := response.ErrorDetail{
		Reason: "internal error",
	}

	data, err := json.Marshal(detail)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if _, exists := decoded["field"]; exists {
		t.Error("field should be omitted when empty (omitempty)")
	}
}

func TestErrorDetail_OmitEmptyReason(t *testing.T) {
	detail := response.ErrorDetail{
		Field: "orderId",
	}

	data, err := json.Marshal(detail)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if _, exists := decoded["reason"]; exists {
		t.Error("reason should be omitted when empty (omitempty)")
	}
}

// ────────────────────────────────────────────────────────────────────────────
// ApiResponse JSON shape
// ────────────────────────────────────────────────────────────────────────────

func TestApiResponse_SuccessShape(t *testing.T) {
	resp := &response.ApiResponse{
		Success:   true,
		Code:      "PARTNER_ORDER_ACCEPTED",
		Message:   "Partner accepted order",
		RequestID: "req-001",
		Data:      map[string]string{"partnerOrderId": "P-abc123"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	requiredFields := []string{"success", "code", "message", "request_id", "data"}
	for _, field := range requiredFields {
		if _, exists := decoded[field]; !exists {
			t.Errorf("expected field %q in ApiResponse JSON", field)
		}
	}
}

func TestApiResponse_ErrorShape(t *testing.T) {
	resp := &response.ApiResponse{
		Success:   false,
		Code:      "INVALID_REQUEST",
		Message:   "Invalid request payload",
		RequestID: "req-002",
		Details:   []response.ErrorDetail{{Field: "orderId", Reason: "required"}},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(data, &decoded)

	if decoded["success"] != false {
		t.Errorf("success = %v, want false", decoded["success"])
	}
	if _, exists := decoded["details"]; !exists {
		t.Error("expected details field in error ApiResponse JSON")
	}
	if _, exists := decoded["data"]; exists {
		t.Error("data should be omitted when nil")
	}
}
