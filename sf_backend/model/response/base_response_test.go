package response

import (
	"net/http"
	"testing"
)

func TestBuildSuccessResponse(t *testing.T) {
	reqID := "req-123"
	data := map[string]string{"key": "value"}
	
	resp := BuildSuccessResponse(OrderCreated, reqID, data)

	if !resp.Success {
		t.Error("expected Success to be true")
	}
	if resp.Code != OrderCreated.Code {
		t.Errorf("expected Code %s, got %s", OrderCreated.Code, resp.Code)
	}
	if resp.Message != OrderCreated.Message {
		t.Errorf("expected Message %s, got %s", OrderCreated.Message, resp.Message)
	}
	if resp.RequestID != reqID {
		t.Errorf("expected RequestID %s, got %s", reqID, resp.RequestID)
	}
	if resp.HTTPStatusCode != OrderCreated.HttpStatusCode {
		t.Errorf("expected HTTPStatusCode %d, got %d", OrderCreated.HttpStatusCode, resp.HTTPStatusCode)
	}
	
	respData := resp.Data.(map[string]string)
	if respData["key"] != "value" {
		t.Errorf("expected data key 'value', got %s", respData["key"])
	}
}

func TestBuildErrorResponse(t *testing.T) {
	reqID := "req-456"
	details := []ErrorDetail{{Field: "orderId", Reason: "required"}}
	
	resp := BuildErrorResponse(InvalidRequest, reqID, details)

	if resp.Success {
		t.Error("expected Success to be false")
	}
	if resp.Code != InvalidRequest.Code {
		t.Errorf("expected Code %s, got %s", InvalidRequest.Code, resp.Code)
	}
	if resp.Message != InvalidRequest.Message {
		t.Errorf("expected Message %s, got %s", InvalidRequest.Message, resp.Message)
	}
	if resp.RequestID != reqID {
		t.Errorf("expected RequestID %s, got %s", reqID, resp.RequestID)
	}
	if resp.HTTPStatusCode != InvalidRequest.HttpStatusCode {
		t.Errorf("expected HTTPStatusCode %d, got %d", InvalidRequest.HttpStatusCode, resp.HTTPStatusCode)
	}
	
	respDetails := resp.Details.([]ErrorDetail)
	if len(respDetails) != 1 || respDetails[0].Field != "orderId" {
		t.Error("details mismatch")
	}
}

func TestResponseCodes(t *testing.T) {
	if OrderCreated.HttpStatusCode != http.StatusCreated {
		t.Errorf("OrderCreated should be 201, got %d", OrderCreated.HttpStatusCode)
	}
	if InternalError.HttpStatusCode != http.StatusInternalServerError {
		t.Errorf("InternalError should be 500, got %d", InternalError.HttpStatusCode)
	}
}
