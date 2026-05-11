package repository

import (
	"context"

	"github.com/xlsmart-api/sf-payment/model/aggregate"
)

type RepositoryInterface interface {
	// Payment transaction operations
	CreatePaymentTransaction(ctx context.Context, payment *aggregate.PaymentTransaction) error
	GetPaymentByReferenceID(ctx context.Context, referenceID string) (*aggregate.PaymentTransaction, error)
	UpdatePaymentTransaction(ctx context.Context, payment *aggregate.PaymentTransaction) error

	// Payment webhook operations
	CreatePaymentWebhook(ctx context.Context, webhook *aggregate.PaymentWebhook) error

	// Idempotency operations
	CheckIdempotency(ctx context.Context, key string) (bool, string, error)
	SaveIdempotency(ctx context.Context, key string, endpoint string, statusCode int, responseBody string) error
}
