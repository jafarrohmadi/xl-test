package aggregate

import (
	"time"

	"github.com/google/uuid"
)

type PaymentTransaction struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ReferenceID   string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"reference_id"`
	TransactionID string     `gorm:"type:varchar(100);uniqueIndex" json:"transaction_id"`
	Amount        float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	Status        string     `gorm:"type:varchar(50);not null" json:"status"`
	PaymentMethod string     `gorm:"type:varchar(50)" json:"payment_method"`
	Provider      string     `gorm:"type:varchar(100)" json:"provider"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type PaymentWebhook struct {
	ID                   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PaymentTransactionID uuid.UUID `gorm:"type:uuid;not null" json:"payment_transaction_id"`
	ReferenceID          string    `gorm:"type:varchar(100);not null" json:"reference_id"`
	TransactionID        string    `gorm:"type:varchar(100);not null" json:"transaction_id"`
	Status               string    `gorm:"type:varchar(50);not null" json:"status"`
	Amount               float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	RawPayload           string    `gorm:"type:jsonb" json:"raw_payload"`
	Signature            string    `gorm:"type:text" json:"signature"`
	WebhookTimestamp     time.Time `json:"webhook_timestamp"`
	ProcessedAt          *time.Time `json:"processed_at"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type IdempotencyRecord struct {
	IdempotencyKey string    `gorm:"type:varchar(100);primary_key" json:"idempotency_key"`
	Endpoint       string    `gorm:"type:varchar(255);not null" json:"endpoint"`
	StatusCode     int       `gorm:"not null" json:"status_code"`
	ResponseBody   string    `gorm:"type:jsonb" json:"response_body"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}
