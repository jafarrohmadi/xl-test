package repository

import (
	"context"

	"github.com/xlsmart-api/sf-backend/model/aggregate"
)

type RepositoryInterface interface {
	// Order operations
	CreateOrder(ctx context.Context, order *aggregate.Order) error
	GetOrderByID(ctx context.Context, orderID string) (*aggregate.Order, error)
	GetOrderByReferenceID(ctx context.Context, referenceID string) (*aggregate.Order, error)
	UpdateOrder(ctx context.Context, order *aggregate.Order) error

	// Order items operations
	CreateOrderItems(ctx context.Context, items []aggregate.OrderItem) error

	// Fulfillment operations
	CreateFulfillment(ctx context.Context, fulfillment *aggregate.Fulfillment) error
	GetFulfillmentByReferenceID(ctx context.Context, referenceID string) (*aggregate.Fulfillment, error)
	UpdateFulfillment(ctx context.Context, fulfillment *aggregate.Fulfillment) error

	// Notification operations
	CreateNotification(ctx context.Context, notification *aggregate.Notification) error

	// Idempotency operations
	CheckIdempotency(ctx context.Context, key string) (bool, string, error)
	SaveIdempotency(ctx context.Context, key string, endpoint string, statusCode int, responseBody string) error
}
