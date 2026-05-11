package request

type PaymentRequest struct {
	ReferenceID string  `json:"referenceId" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

type PaymentWebhookRequest struct {
	ReferenceID   string  `json:"referenceId" validate:"required"`
	TransactionID string  `json:"transactionId" validate:"required"`
	Status        string  `json:"status" validate:"required,oneof=SUCCESS FAILED EXPIRED"`
	Amount        float64 `json:"amount" validate:"required,min=0"`
	PaidAt        string  `json:"paidAt"`
}
