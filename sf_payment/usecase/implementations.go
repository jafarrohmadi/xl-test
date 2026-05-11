package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xlsmart-api/sf-payment/model/aggregate"
	"github.com/xlsmart-api/sf-payment/model/request"
)

func (u *UseCase) ProcessPayment(ctx context.Context, req request.PaymentRequest, requestID string, idempotencyKey string) (map[string]interface{}, error) {
	// Check idempotency
	exists, cachedResponse, err := u.Repository.CheckIdempotency(ctx, idempotencyKey)
	if err != nil {
		return nil, err
	}
	if exists {
		var response map[string]interface{}
		json.Unmarshal([]byte(cachedResponse), &response)
		return response, nil
	}

	// Create payment transaction record
	now := time.Now()
	payment := aggregate.PaymentTransaction{
		ReferenceID: req.ReferenceID,
		Amount:      req.Price,
		Status:      "PROCESSING",
		Provider:    "payment-gateway",
	}

	if err := u.Repository.CreatePaymentTransaction(ctx, &payment); err != nil {
		return nil, err
	}

	// TODO: Call payment gateway API to process payment
	// Simulate payment processing
	payment.TransactionID = fmt.Sprintf("TXN-PG-%d", time.Now().Unix())
	payment.Status = "SUCCESS"
	payment.PaidAt = &now

	if err := u.Repository.UpdatePaymentTransaction(ctx, &payment); err != nil {
		return nil, err
	}

	// Build response
	response := map[string]interface{}{
		"referenceId":   payment.ReferenceID,
		"paymentStatus": payment.Status,
		"paidAt":        payment.PaidAt.Format(time.RFC3339),
	}

	// Save idempotency
	responseJSON, _ := json.Marshal(response)
	u.Repository.SaveIdempotency(ctx, idempotencyKey, "/sf-payment/payments/request", 200, string(responseJSON))

	return response, nil
}

func (u *UseCase) ProcessWebhook(ctx context.Context, req request.PaymentWebhookRequest, requestID string, idempotencyKey string, signature string, timestamp string) error {
	// Check idempotency
	exists, _, err := u.Repository.CheckIdempotency(ctx, idempotencyKey)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("already processed")
	}

	// TODO: Validate HMAC signature
	// TODO: Verify timestamp (replay window < 5 min)

	// Find payment transaction by reference ID
	payment, err := u.Repository.GetPaymentByReferenceID(ctx, req.ReferenceID)
	if err != nil {
		return fmt.Errorf("payment not found")
	}

	// Create webhook record
	rawPayloadJSON, _ := json.Marshal(req)
	webhookTimestamp, _ := time.Parse(time.RFC3339, timestamp)
	now := time.Now()
	webhook := aggregate.PaymentWebhook{
		PaymentTransactionID: payment.ID,
		ReferenceID:          req.ReferenceID,
		TransactionID:        req.TransactionID,
		Status:               req.Status,
		Amount:               req.Amount,
		RawPayload:           string(rawPayloadJSON),
		Signature:            signature,
		WebhookTimestamp:     webhookTimestamp,
		ProcessedAt:          &now,
	}

	if err := u.Repository.CreatePaymentWebhook(ctx, &webhook); err != nil {
		return err
	}

	// Update payment transaction status
	payment.Status = req.Status
	if req.PaidAt != "" {
		paidAt, _ := time.Parse(time.RFC3339, req.PaidAt)
		payment.PaidAt = &paidAt
	}

	if err := u.Repository.UpdatePaymentTransaction(ctx, payment); err != nil {
		return err
	}

	// Save idempotency
	u.Repository.SaveIdempotency(ctx, idempotencyKey, "/sf-payment/payments/webhook", 200, "{}")

	// TODO: Notify SF Backend of payment status update (if status changed)

	return nil
}
