package exceptions

import (
	"context"
	"errors"
	"strings"

	"github.com/xlsmart-api/partner/model/response"
	"github.com/go-playground/validator/v10"
)

var (
	OrderNotFoundError         = errors.New("order not found")
	PaymentNotFoundError       = errors.New("payment not found")
	DuplicateOrderError        = errors.New("duplicate order")
	DuplicatePaymentError      = errors.New("duplicate payment")
	AlreadyProcessedError      = errors.New("already processed")
	InvalidSignatureError      = errors.New("invalid signature")
	UnauthorizedError          = errors.New("unauthorized")
	PartnerSubmitFailedError   = errors.New("partner submit failed")
	PaymentFailedError         = errors.New("payment failed")
	RateLimitedError           = errors.New("rate limited")
)

func HandleError(ctx context.Context, err error, requestID string) *response.ApiResponse {
	// Handle validation errors
	if _, ok := err.(validator.ValidationErrors); ok {
		return response.BuildErrorResponse(response.InvalidRequest, requestID, nil)
	}

	// Handle specific errors
	if strings.Contains(err.Error(), OrderNotFoundError.Error()) {
		return response.BuildErrorResponse(response.OrderNotFound, requestID, nil)
	}

	if strings.Contains(err.Error(), PaymentNotFoundError.Error()) {
		return response.BuildErrorResponse(response.PaymentNotFound, requestID, nil)
	}

	if strings.Contains(err.Error(), DuplicateOrderError.Error()) {
		return response.BuildErrorResponse(response.DuplicateOrder, requestID, nil)
	}

	if strings.Contains(err.Error(), DuplicatePaymentError.Error()) {
		return response.BuildErrorResponse(response.DuplicatePaymentRequest, requestID, nil)
	}

	if strings.Contains(err.Error(), AlreadyProcessedError.Error()) {
		return response.BuildErrorResponse(response.AlreadyProcessed, requestID, nil)
	}

	if strings.Contains(err.Error(), InvalidSignatureError.Error()) {
		return response.BuildErrorResponse(response.InvalidSignature, requestID, nil)
	}

	if strings.Contains(err.Error(), UnauthorizedError.Error()) {
		return response.BuildErrorResponse(response.Unauthorized, requestID, nil)
	}

	if strings.Contains(err.Error(), PartnerSubmitFailedError.Error()) {
		return response.BuildErrorResponse(response.PartnerSubmitFailed, requestID, nil)
	}

	if strings.Contains(err.Error(), PaymentFailedError.Error()) {
		return response.BuildErrorResponse(response.PaymentFailed, requestID, nil)
	}

	if strings.Contains(err.Error(), RateLimitedError.Error()) {
		return response.BuildErrorResponse(response.RateLimited, requestID, nil)
	}

	// Handle other errors as internal server errors
	return response.BuildErrorResponse(response.InternalError, requestID, nil)
}
