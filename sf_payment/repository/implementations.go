package repository

import (
	"context"
	"errors"
	"time"

	"github.com/xlsmart-api/sf-payment/model/aggregate"
	"gorm.io/gorm"
)

func (r *Repository) CreatePaymentTransaction(ctx context.Context, payment *aggregate.PaymentTransaction) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *Repository) GetPaymentByReferenceID(ctx context.Context, referenceID string) (*aggregate.PaymentTransaction, error) {
	var payment aggregate.PaymentTransaction
	err := r.db.WithContext(ctx).Where("reference_id = ?", referenceID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *Repository) UpdatePaymentTransaction(ctx context.Context, payment *aggregate.PaymentTransaction) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *Repository) CreatePaymentWebhook(ctx context.Context, webhook *aggregate.PaymentWebhook) error {
	return r.db.WithContext(ctx).Create(webhook).Error
}

func (r *Repository) CheckIdempotency(ctx context.Context, key string) (bool, string, error) {
	var record aggregate.IdempotencyRecord
	err := r.db.WithContext(ctx).Where("idempotency_key = ? AND expires_at > ?", key, time.Now()).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", nil
		}
		return false, "", err
	}
	return true, record.ResponseBody, nil
}

func (r *Repository) SaveIdempotency(ctx context.Context, key string, endpoint string, statusCode int, responseBody string) error {
	record := aggregate.IdempotencyRecord{
		IdempotencyKey: key,
		Endpoint:       endpoint,
		StatusCode:     statusCode,
		ResponseBody:   responseBody,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}
	return r.db.WithContext(ctx).Create(&record).Error
}
