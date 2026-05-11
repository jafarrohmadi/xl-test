package request

type PartnerGoods struct {
	SKU string `json:"sku" validate:"required"`
	Qty int    `json:"qty" validate:"required,min=1"`
}

type PartnerSubmitRequest struct {
	ReferenceID string         `json:"referenceId" validate:"required"`
	OrderID     string         `json:"orderId" validate:"required"`
	Goods       []PartnerGoods `json:"goods" validate:"required,min=1,dive"`
	TotalPrice  float64        `json:"totalPrice" validate:"required,min=0"`
}

type FulfillmentRequest struct {
	ReferenceID    string `json:"referenceId" validate:"required"`
	PartnerOrderID string `json:"partnerOrderId" validate:"required"`
}
