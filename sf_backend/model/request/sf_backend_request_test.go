package request

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestSubmitOrderRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     SubmitOrderRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: SubmitOrderRequest{
				OrderID:   "order-123",
				PartnerID: "partner-456",
				Goods: []GoodsItem{
					{SKU: "sku-1", Name: "Item 1", Qty: 1},
				},
				TotalPrice: 100.0,
			},
			wantErr: false,
		},
		{
			name: "missing order id",
			req: SubmitOrderRequest{
				PartnerID: "partner-456",
				Goods: []GoodsItem{
					{SKU: "sku-1", Name: "Item 1", Qty: 1},
				},
				TotalPrice: 100.0,
			},
			wantErr: true,
		},
		{
			name: "empty goods list",
			req: SubmitOrderRequest{
				OrderID:   "order-123",
				PartnerID: "partner-456",
				Goods:     []GoodsItem{},
				TotalPrice: 100.0,
			},
			wantErr: true,
		},
		{
			name: "invalid quantity",
			req: SubmitOrderRequest{
				OrderID:   "order-123",
				PartnerID: "partner-456",
				Goods: []GoodsItem{
					{SKU: "sku-1", Name: "Item 1", Qty: 0},
				},
				TotalPrice: 100.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmitOrderRequest validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFulfillmentCallbackRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     FulfillmentCallbackRequest
		wantErr bool
	}{
		{
			name: "valid success callback",
			req: FulfillmentCallbackRequest{
				ReferenceID:    "ref-123",
				PartnerOrderID: "p-order-123",
				Status:         "SUCCESS",
			},
			wantErr: false,
		},
		{
			name: "valid failed callback",
			req: FulfillmentCallbackRequest{
				ReferenceID:    "ref-123",
				PartnerOrderID: "p-order-123",
				Status:         "FAILED",
				FailureReason:  "Out of stock",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			req: FulfillmentCallbackRequest{
				ReferenceID:    "ref-123",
				PartnerOrderID: "p-order-123",
				Status:         "PENDING",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FulfillmentCallbackRequest validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotificationEvent_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     NotificationEvent
		wantErr bool
	}{
		{
			name: "valid notification event",
			req: NotificationEvent{
				EventType:   "ORDER_PAID",
				ReferenceID: "ref-123",
				UserID:      "user-123",
				Channels:    []string{"EMAIL", "SMS"},
			},
			wantErr: false,
		},
		{
			name: "empty channels",
			req: NotificationEvent{
				EventType:   "ORDER_PAID",
				ReferenceID: "ref-123",
				UserID:      "user-123",
				Channels:    []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationEvent validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
