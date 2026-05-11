package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xlsmart-api/sf-backend/model/aggregate"
	"github.com/xlsmart-api/sf-backend/model/request"
)

func (u *UseCase) SubmitOrder(ctx context.Context, req request.SubmitOrderRequest, requestID string, idempotencyKey string) (map[string]interface{}, error) {
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

	// Generate reference ID
	referenceID := fmt.Sprintf("REF-%s", uuid.New().String()[:8])

	// Create order record
	order := aggregate.Order{
		OrderID:       req.OrderID,
		ReferenceID:   referenceID,
		UserID:        "user-123", // TODO: Extract from JWT token
		PartnerID:     req.PartnerID,
		TotalPrice:    req.TotalPrice,
		Status:        "SUBMITTED",
		PaymentStatus: "PENDING",
	}

	if err := u.Repository.CreateOrder(ctx, &order); err != nil {
		return nil, err
	}

	// Create order items
	var items []aggregate.OrderItem
	for _, goods := range req.Goods {
		items = append(items, aggregate.OrderItem{
			OrderID:     order.ID,
			SKU:         goods.SKU,
			Name:        goods.Name,
			Description: goods.Desc,
			Quantity:    goods.Qty,
		})
	}
	if err := u.Repository.CreateOrderItems(ctx, items); err != nil {
		return nil, err
	}

	partnerURL := os.Getenv("PARTNER_URL")
	if partnerURL == "" {
		partnerURL = "http://localhost:8083"
	}
	paymentURL := os.Getenv("SF_PAYMENT_URL")
	if paymentURL == "" {
		paymentURL = "http://localhost:8082"
	}

	// 1. Call Partner module to submit order
	partnerReqBody := map[string]interface{}{
		"referenceId": referenceID,
		"orderId":     req.OrderID,
		"goods":       req.Goods,
		"totalPrice":  req.TotalPrice,
	}
	partnerResData, err := sendHTTPRequest(ctx, "POST", partnerURL+"/partners/orders", partnerReqBody, requestID, idempotencyKey+"-partner-submit")
	if err != nil {
		return nil, fmt.Errorf("failed to submit partner order: %w", err)
	}
	partnerOrderID := partnerResData["partnerOrderId"].(string)

	// Update order with partner order ID
	order.PartnerOrderID = partnerOrderID
	order.Status = "PARTNER_ACCEPTED"
	if err := u.Repository.UpdateOrder(ctx, &order); err != nil {
		return nil, err
	}

	// 2. Call SF Payment module to process payment
	paymentReqBody := map[string]interface{}{
		"referenceId": referenceID,
		"price":       req.TotalPrice,
	}
	paymentResData, err := sendHTTPRequest(ctx, "POST", paymentURL+"/payments", paymentReqBody, requestID, idempotencyKey+"-payment")
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}
	paymentStatus := paymentResData["paymentStatus"].(string)

	// Update order payment status
	order.PaymentStatus = paymentStatus
	if err := u.Repository.UpdateOrder(ctx, &order); err != nil {
		return nil, err
	}

	// 3. Call Partner module to request fulfillment
	fulfillmentReqBody := map[string]interface{}{
		"referenceId":    referenceID,
		"partnerOrderId": partnerOrderID,
	}
	_, err = sendHTTPRequest(ctx, "POST", partnerURL+"/partners/fulfillment", fulfillmentReqBody, requestID, idempotencyKey+"-fulfillment")
	if err != nil {
		return nil, fmt.Errorf("failed to request fulfillment: %w", err)
	}

	// Build response
	response := map[string]interface{}{
		"orderId":        order.OrderID,
		"referenceId":    order.ReferenceID,
		"partnerOrderId": order.PartnerOrderID,
		"paymentStatus":  order.PaymentStatus,
	}

	// Save idempotency
	responseJSON, _ := json.Marshal(response)
	u.Repository.SaveIdempotency(ctx, idempotencyKey, "/sf-backend/orders/submit", 201, string(responseJSON))

	return response, nil
}

func sendHTTPRequest(ctx context.Context, method, url string, body interface{}, requestID, idempotencyKey string) (map[string]interface{}, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if requestID != "" {
		req.Header.Set("X-Request-Id", requestID)
	}
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	var apiResp struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

func (u *UseCase) ProcessFulfillmentCallback(ctx context.Context, req request.FulfillmentCallbackRequest, requestID string, idempotencyKey string, signature string, timestamp string) error {
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

	// Find order by reference ID
	order, err := u.Repository.GetOrderByReferenceID(ctx, req.ReferenceID)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	// Create or update fulfillment record
	voucherJSON, _ := json.Marshal(req.Voucher)
	now := time.Now()
	fulfillment := aggregate.Fulfillment{
		OrderID:        order.ID,
		ReferenceID:    req.ReferenceID,
		PartnerOrderID: req.PartnerOrderID,
		Status:         req.Status,
		VoucherData:    string(voucherJSON),
		FailureReason:  req.FailureReason,
		FulfilledAt:    &now,
	}

	if err := u.Repository.CreateFulfillment(ctx, &fulfillment); err != nil {
		return err
	}

	// Update order status
	order.Status = "FULFILLMENT_SUCCESS"
	if err := u.Repository.UpdateOrder(ctx, order); err != nil {
		return err
	}

	// Save idempotency
	u.Repository.SaveIdempotency(ctx, idempotencyKey, "/sf-backend/fulfillment/callback", 200, "{}")

	// Trigger notification
	notificationReq := request.NotificationEvent{
		EventType:   "fulfillment.completed",
		ReferenceID: req.ReferenceID,
		UserID:      order.UserID,
		Channels:    []string{"email", "push"},
		Data: map[string]interface{}{
			"order_id":     order.OrderID,
			"voucher_code": req.Voucher.Code,
		},
	}
	u.TriggerNotification(ctx, notificationReq, requestID)

	return nil
}

func (u *UseCase) TriggerNotification(ctx context.Context, req request.NotificationEvent, requestID string) error {
	// Find order by reference ID
	order, err := u.Repository.GetOrderByReferenceID(ctx, req.ReferenceID)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	// Create notification record
	dataJSON, _ := json.Marshal(req.Data)
	notification := aggregate.Notification{
		OrderID:     order.ID,
		ReferenceID: req.ReferenceID,
		UserID:      req.UserID,
		EventType:   req.EventType,
		Channels:    strings.Join(req.Channels, ","),
		Data:        string(dataJSON),
		Status:      "QUEUED",
	}

	if err := u.Repository.CreateNotification(ctx, &notification); err != nil {
		return err
	}

	// TODO: Enqueue notification job to message queue

	return nil
}
