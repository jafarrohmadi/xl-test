package usecase

import (
	"context"

	"github.com/xlsmart-api/partner/model/request"
)

type UseCaseInterface interface {
	// Submit order to Partner
	SubmitOrder(ctx context.Context, req request.PartnerSubmitRequest, requestID string, idempotencyKey string) (map[string]interface{}, error)

	// Process fulfillment request
	ProcessFulfillment(ctx context.Context, req request.FulfillmentRequest, requestID string, idempotencyKey string) error
}
