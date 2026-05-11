package response

import (
	"net/http"
)

type ApiResponse struct {
	Success   bool        `json:"success"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
	Data           interface{} `json:"data,omitempty"`
	Details        interface{} `json:"details,omitempty"`
	HTTPStatusCode int         `json:"-"`
}

type ErrorDetail struct {
	Field  string `json:"field,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type ResponseCode struct {
	HttpStatusCode int
	Code           string
	Message        string
}

// Success codes
var OrderCreated = ResponseCode{http.StatusCreated, "ORDER_CREATED", "Order created successfully"}
var FulfillmentCallbackAccepted = ResponseCode{http.StatusOK, "FULFILLMENT_CALLBACK_ACCEPTED", "Callback processed"}
var NotificationQueued = ResponseCode{http.StatusAccepted, "NOTIFICATION_QUEUED", "Notification job queued"}
var PartnerOrderAccepted = ResponseCode{http.StatusOK, "PARTNER_ORDER_ACCEPTED", "Partner accepted order"}
var FulfillmentInProgress = ResponseCode{http.StatusOK, "FULFILLMENT_IN_PROGRESS", "Fulfillment process started"}

// Error codes
var InvalidRequest = ResponseCode{http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload"}
var InvalidSignature = ResponseCode{http.StatusBadRequest, "INVALID_SIGNATURE", "Invalid HMAC signature"}
var Unauthorized = ResponseCode{http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or missing authentication"}
var NotFound = ResponseCode{http.StatusNotFound, "NOT_FOUND", "Resource not found"}
var OrderNotFound = ResponseCode{http.StatusNotFound, "ORDER_NOT_FOUND", "Order not found"}
var PaymentNotFound = ResponseCode{http.StatusNotFound, "PAYMENT_NOT_FOUND", "Payment transaction not found"}
var DuplicateOrder = ResponseCode{http.StatusConflict, "DUPLICATE_ORDER", "Order already exists"}
var AlreadyProcessed = ResponseCode{http.StatusConflict, "ALREADY_PROCESSED", "Request already processed"}
var DuplicatePaymentRequest = ResponseCode{http.StatusConflict, "DUPLICATE_PAYMENT_REQUEST", "Payment already processed"}
var PartnerSubmitFailed = ResponseCode{http.StatusUnprocessableEntity, "PARTNER_SUBMIT_FAILED", "Partner rejected submit order"}
var PaymentFailed = ResponseCode{http.StatusPaymentRequired, "PAYMENT_FAILED", "Payment rejected by provider"}
var RateLimited = ResponseCode{http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests"}
var InternalError = ResponseCode{http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"}

func BuildSuccessResponse(responseCode ResponseCode, requestID string, data interface{}) *ApiResponse {
	return &ApiResponse{
		Success:   true,
		Code:      responseCode.Code,
		Message:        responseCode.Message,
		RequestID:      requestID,
		Data:           data,
		HTTPStatusCode: responseCode.HttpStatusCode,
	}
}

func BuildErrorResponse(responseCode ResponseCode, requestID string, details interface{}) *ApiResponse {
	return &ApiResponse{
		Success:   false,
		Code:      responseCode.Code,
		Message:        responseCode.Message,
		RequestID:      requestID,
		Details:        details,
		HTTPStatusCode: responseCode.HttpStatusCode,
	}
}
