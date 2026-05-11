package request

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestPaymentRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     PaymentRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: PaymentRequest{
				ReferenceID: "ref-123",
				Price:       100.0,
			},
			wantErr: false,
		},
		{
			name: "missing reference id",
			req: PaymentRequest{
				Price: 100.0,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			req: PaymentRequest{
				ReferenceID: "ref-123",
				Price:       -10.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentRequest validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaymentWebhookRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     PaymentWebhookRequest
		wantErr bool
	}{
		{
			name: "valid success webhook",
			req: PaymentWebhookRequest{
				ReferenceID:   "ref-123",
				TransactionID: "trx-456",
				Status:        "SUCCESS",
				Amount:        100.0,
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			req: PaymentWebhookRequest{
				ReferenceID:   "ref-123",
				TransactionID: "trx-456",
				Status:        "PENDING",
				Amount:        100.0,
			},
			wantErr: true,
		},
		{
			name: "missing transaction id",
			req: PaymentWebhookRequest{
				ReferenceID: "ref-123",
				Status:      "SUCCESS",
				Amount:      100.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentWebhookRequest validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
