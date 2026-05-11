package aggregate

import (
	"time"

	"github.com/google/uuid"
)

func (Order) TableName() string {
	return "partner_orders"
}

type Order struct {
	ID              uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ReferenceID     string      `gorm:"type:varchar(100);uniqueIndex;not null" json:"reference_id"`
	OrderID         string      `gorm:"type:varchar(100);not null" json:"order_id"`
	PartnerOrderID  string      `gorm:"type:varchar(100);uniqueIndex;not null" json:"partner_order_id"`
	PartnerID       string      `gorm:"type:varchar(100);not null" json:"partner_id"`
	TotalPrice      float64     `gorm:"type:decimal(15,2);not null" json:"total_price"`
	Status          string      `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	SubmittedAt     *time.Time  `gorm:"autoCreateTime" json:"submitted_at"`
	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	OrderItems      []OrderItem `gorm:"foreignKey:PartnerOrderID;references:ID" json:"order_items,omitempty"`
}

func (OrderItem) TableName() string {
	return "partner_order_items"
}

type OrderItem struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PartnerOrderID uuid.UUID `gorm:"type:uuid;not null;column:partner_order_id" json:"partner_order_id"`
	SKU            string    `gorm:"type:varchar(100);not null" json:"sku"`
	Quantity       int       `gorm:"not null" json:"quantity"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Fulfillment) TableName() string {
	return "partner_fulfillments"
}

type Fulfillment struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PartnerOrderID      uuid.UUID  `gorm:"type:uuid;not null;column:partner_order_id" json:"partner_order_id"`
	ReferenceID         string     `gorm:"type:varchar(100);not null" json:"reference_id"`
	Status              string     `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	VoucherCode         string     `gorm:"type:varchar(255)" json:"voucher_code"`
	VoucherSerialNumber string     `gorm:"type:varchar(255)" json:"voucher_serial_number"`
	FailureReason       string     `gorm:"type:text" json:"failure_reason"`
	RequestedAt         time.Time  `gorm:"autoCreateTime" json:"requested_at"`
	FulfilledAt         *time.Time `json:"fulfilled_at"`
	CallbackSentAt      *time.Time `json:"callback_sent_at"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
