package usecase

import (
	"context"

	"github.com/xlsmart-api/sf-payment/model/request"
)

type UseCaseInterface interface {
	// Process payment request from SF Backend
	ProcessPayment(ctx context.Context, req request.PaymentRequest, requestID string, idempotencyKey string) (map[string]interface{}, error)

	// Process payment gateway webhook
	ProcessWebhook(ctx context.Context, req request.PaymentWebhookRequest, requestID string, idempotencyKey string, signature string, timestamp string) error
}
