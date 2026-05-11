package usecase

import (
	"context"

	"github.com/xlsmart-api/sf-backend/model/request"
)

type UseCaseInterface interface {
	// Submit order orchestration
	SubmitOrder(ctx context.Context, req request.SubmitOrderRequest, requestID string, idempotencyKey string) (map[string]interface{}, error)

	// Process fulfillment callback from Partner
	ProcessFulfillmentCallback(ctx context.Context, req request.FulfillmentCallbackRequest, requestID string, idempotencyKey string, signature string, timestamp string) error

	// Trigger notification
	TriggerNotification(ctx context.Context, req request.NotificationEvent, requestID string) error
}
