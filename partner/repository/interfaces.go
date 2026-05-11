package repository

import (
	"context"

	"github.com/xlsmart-api/partner/model/aggregate"
)

type RepositoryInterface interface {
	// Partner order operations
	CreatePartnerOrder(ctx context.Context, order *aggregate.Order) error
	GetPartnerOrderByReferenceID(ctx context.Context, referenceID string) (*aggregate.Order, error)
	UpdatePartnerOrder(ctx context.Context, order *aggregate.Order) error

	// Fulfillment operations
	CreateFulfillment(ctx context.Context, fulfillment *aggregate.Fulfillment) error
	GetFulfillmentByPartnerOrderID(ctx context.Context, partnerOrderID string) (*aggregate.Fulfillment, error)
	UpdateFulfillment(ctx context.Context, fulfillment *aggregate.Fulfillment) error

	// Idempotency operations
	CheckIdempotency(ctx context.Context, key string) (bool, string, error)
	SaveIdempotency(ctx context.Context, key string, endpoint string, statusCode int, responseBody string) error
}
