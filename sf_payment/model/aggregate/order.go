package aggregate

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID         string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"order_id"`
	ReferenceID     string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"reference_id"`
	UserID          string     `gorm:"type:varchar(100);not null" json:"user_id"`
	PartnerID       string     `gorm:"type:varchar(100);not null" json:"partner_id"`
	PartnerOrderID  string     `gorm:"type:varchar(100)" json:"partner_order_id"`
	TotalPrice      float64    `gorm:"type:decimal(15,2);not null" json:"total_price"`
	Status          string     `gorm:"type:varchar(50);not null" json:"status"`
	PaymentStatus   string     `gorm:"type:varchar(50);not null" json:"payment_status"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	OrderItems      []OrderItem `gorm:"foreignKey:OrderID;references:ID" json:"order_items,omitempty"`
}

type OrderItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	SKU         string    `gorm:"type:varchar(100);not null" json:"sku"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Fulfillment struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID         uuid.UUID  `gorm:"type:uuid;not null" json:"order_id"`
	ReferenceID     string     `gorm:"type:varchar(100);not null" json:"reference_id"`
	PartnerOrderID  string     `gorm:"type:varchar(100);not null" json:"partner_order_id"`
	Status          string     `gorm:"type:varchar(50);not null" json:"status"`
	VoucherData     string     `gorm:"type:jsonb" json:"voucher_data"`
	FailureReason   string     `gorm:"type:text" json:"failure_reason"`
	FulfilledAt     *time.Time `json:"fulfilled_at"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Notification struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID     uuid.UUID  `gorm:"type:uuid;not null" json:"order_id"`
	ReferenceID string     `gorm:"type:varchar(100);not null" json:"reference_id"`
	UserID      string     `gorm:"type:varchar(100);not null" json:"user_id"`
	EventType   string     `gorm:"type:varchar(100);not null" json:"event_type"`
	Channels    string     `gorm:"type:varchar(255);not null" json:"channels"`
	Data        string     `gorm:"type:jsonb" json:"data"`
	Status      string     `gorm:"type:varchar(50);not null" json:"status"`
	SentAt      *time.Time `json:"sent_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
