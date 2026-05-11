package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/xlsmart-api/partner/model/aggregate"
	"github.com/xlsmart-api/partner/model/request"
)

func (u *UseCase) SubmitOrder(ctx context.Context, req request.PartnerSubmitRequest, requestID string, idempotencyKey string) (map[string]interface{}, error) {
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

	// Generate partner order ID
	partnerOrderID := fmt.Sprintf("P-%s", uuid.New().String()[:6])

	// Create partner order record
	order := aggregate.Order{
		OrderID:        req.OrderID,
		ReferenceID:    req.ReferenceID,
		PartnerOrderID: partnerOrderID,
		PartnerID:      "DEFAULT_PARTNER",
		TotalPrice:     req.TotalPrice,
		Status:         "ACCEPTED",
	}

	if err := u.Repository.CreatePartnerOrder(ctx, &order); err != nil {
		return nil, err
	}

	// Build response
	response := map[string]interface{}{
		"partnerOrderId": partnerOrderID,
		"status":         "ACCEPTED",
	}

	// Save idempotency
	responseJSON, _ := json.Marshal(response)
	u.Repository.SaveIdempotency(ctx, idempotencyKey, "/partner/orders/submit", 200, string(responseJSON))

	return response, nil
}

func (u *UseCase) ProcessFulfillment(ctx context.Context, req request.FulfillmentRequest, requestID string, idempotencyKey string) error {
	// Check idempotency
	exists, _, err := u.Repository.CheckIdempotency(ctx, idempotencyKey)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Find partner order by reference ID
	order, err := u.Repository.GetPartnerOrderByReferenceID(ctx, req.ReferenceID)
	if err != nil {
		return fmt.Errorf("partner order not found")
	}

	// Create fulfillment record
	fulfillment := aggregate.Fulfillment{
		PartnerOrderID: order.ID,
		ReferenceID:    req.ReferenceID,
		Status:         "IN_PROGRESS",
	}

	if err := u.Repository.CreateFulfillment(ctx, &fulfillment); err != nil {
		return err
	}

	// Save idempotency
	u.Repository.SaveIdempotency(ctx, idempotencyKey, "/partner/orders/fulfillment", 200, "{}")

	// Async fulfillment processing
	go func(refID, partnerOrderID string) {
		// Simulate processing delay
		time.Sleep(2 * time.Second)

		callbackURL := os.Getenv("SF_BACKEND_CALLBACK_URL")
		if callbackURL == "" {
			callbackURL = "http://localhost:8081/orders/fulfillment/callback"
		}

		callbackReq := map[string]interface{}{
			"referenceId":    refID,
			"partnerOrderId": partnerOrderID,
			"status":         "SUCCESS",
			"voucher": map[string]string{
				"code":         fmt.Sprintf("V-%s", uuid.New().String()[:6]),
				"serialNumber": fmt.Sprintf("SN-%s", uuid.New().String()[:8]),
			},
		}

		jsonBody, _ := json.Marshal(callbackReq)
		req, err := http.NewRequest("POST", callbackURL, bytes.NewBuffer(jsonBody))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Request-Id", uuid.New().String())
			req.Header.Set("Idempotency-Key", uuid.New().String())
			req.Header.Set("X-Signature", "dummy-signature")
			req.Header.Set("X-Signature-Timestamp", time.Now().Format(time.RFC3339))
			client := &http.Client{Timeout: 5 * time.Second}
			client.Do(req)
		}
	}(req.ReferenceID, order.PartnerOrderID)

	return nil
}
