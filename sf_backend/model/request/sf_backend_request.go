package request

type GoodsItem struct {
	SKU  string `json:"sku" validate:"required"`
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc"`
	Qty  int    `json:"qty" validate:"required,min=1"`
}

type SubmitOrderRequest struct {
	OrderID    string      `json:"orderId" validate:"required"`
	PartnerID  string      `json:"partnerId" validate:"required"`
	Goods      []GoodsItem `json:"goods" validate:"required,min=1,dive"`
	TotalPrice float64     `json:"totalPrice" validate:"required,min=0"`
}

type FulfillmentCallbackRequest struct {
	ReferenceID    string        `json:"referenceId" validate:"required"`
	PartnerOrderID string        `json:"partnerOrderId" validate:"required"`
	Status         string        `json:"status" validate:"required,oneof=SUCCESS FAILED"`
	Voucher        *VoucherData  `json:"voucher"`
	FailureReason  string        `json:"failureReason"`
}

type VoucherData struct {
	Code         string `json:"code"`
	SerialNumber string `json:"serialNumber"`
}

type NotificationEvent struct {
	EventType   string                 `json:"event_type" validate:"required"`
	ReferenceID string                 `json:"reference_id" validate:"required"`
	UserID      string                 `json:"user_id" validate:"required"`
	Channels    []string               `json:"channels" validate:"required,min=1"`
	Data        map[string]interface{} `json:"data"`
}
