package repository

import (
	"context"
	"errors"
	"time"

	"github.com/xlsmart-api/partner/model/aggregate"
	"gorm.io/gorm"
)

func (r *Repository) CreatePartnerOrder(ctx context.Context, order *aggregate.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *Repository) GetPartnerOrderByReferenceID(ctx context.Context, referenceID string) (*aggregate.Order, error) {
	var order aggregate.Order
	err := r.db.WithContext(ctx).Where("reference_id = ?", referenceID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) UpdatePartnerOrder(ctx context.Context, order *aggregate.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *Repository) CreateFulfillment(ctx context.Context, fulfillment *aggregate.Fulfillment) error {
	return r.db.WithContext(ctx).Create(fulfillment).Error
}

func (r *Repository) GetFulfillmentByPartnerOrderID(ctx context.Context, partnerOrderID string) (*aggregate.Fulfillment, error) {
	var fulfillment aggregate.Fulfillment
	err := r.db.WithContext(ctx).Where("partner_order_id = ?", partnerOrderID).First(&fulfillment).Error
	if err != nil {
		return nil, err
	}
	return &fulfillment, nil
}

func (r *Repository) UpdateFulfillment(ctx context.Context, fulfillment *aggregate.Fulfillment) error {
	return r.db.WithContext(ctx).Save(fulfillment).Error
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
