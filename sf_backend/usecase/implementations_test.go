package usecase

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xlsmart-api/sf-backend/model/aggregate"
	"github.com/xlsmart-api/sf-backend/model/request"
	repo_mocks "github.com/xlsmart-api/sf-backend/repository/mocks"
)

func TestSubmitOrder(t *testing.T) {
	mockRepo := new(repo_mocks.MockRepository)
	uc := NewUseCase(NewUseCaseOptions{
		Repository: mockRepo,
	})

	ctx := context.Background()
	req := request.SubmitOrderRequest{
		OrderID:    "order-123",
		PartnerID:  "partner-456",
		TotalPrice: 100000,
		Goods: []request.GoodsItem{
			{SKU: "SKU1", Name: "Item 1", Qty: 1},
		},
	}
	requestID := "req-123"
	idempotencyKey := "idem-123"

	t.Run("success", func(t *testing.T) {
		// Mock idempotency check
		mockRepo.On("CheckIdempotency", ctx, idempotencyKey).Return(false, "", nil).Once()

		// Mock order creation
		mockRepo.On("CreateOrder", ctx, mock.AnythingOfType("*aggregate.Order")).Return(nil).Once()
		mockRepo.On("CreateOrderItems", ctx, mock.AnythingOfType("[]aggregate.OrderItem")).Return(nil).Once()

		// Mock HTTP calls using httptest
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			
			var respData map[string]interface{}
			if r.URL.Path == "/partners/orders" {
				respData = map[string]interface{}{"partnerOrderId": "p-order-123"}
			} else if r.URL.Path == "/payments" {
				respData = map[string]interface{}{"paymentStatus": "PAID"}
			} else if r.URL.Path == "/partners/fulfillment" {
				respData = map[string]interface{}{}
			}

			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    respData,
			})
		}))
		defer server.Close()

		// Point env to test server
		os.Setenv("PARTNER_URL", server.URL)
		os.Setenv("SF_PAYMENT_URL", server.URL)
		defer os.Unsetenv("PARTNER_URL")
		defer os.Unsetenv("SF_PAYMENT_URL")

		// Mock updates
		mockRepo.On("UpdateOrder", ctx, mock.AnythingOfType("*aggregate.Order")).Return(nil).Twice()
		mockRepo.On("SaveIdempotency", ctx, idempotencyKey, "/sf-backend/orders/submit", 201, mock.Anything).Return(nil).Once()

		resp, err := uc.SubmitOrder(ctx, req, requestID, idempotencyKey)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "order-123", resp["orderId"])
		assert.Equal(t, "p-order-123", resp["partnerOrderId"])
		assert.Equal(t, "PAID", resp["paymentStatus"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("idempotency exists", func(t *testing.T) {
		cachedResp := `{"orderId":"order-123","status":"SUBMITTED"}`
		mockRepo.On("CheckIdempotency", ctx, idempotencyKey).Return(true, cachedResp, nil).Once()

		resp, err := uc.SubmitOrder(ctx, req, requestID, idempotencyKey)

		assert.NoError(t, err)
		assert.Equal(t, "order-123", resp["orderId"])
		mockRepo.AssertExpectations(t)
	})
}

func TestProcessFulfillmentCallback(t *testing.T) {
	mockRepo := new(repo_mocks.MockRepository)
	uc := NewUseCase(NewUseCaseOptions{
		Repository: mockRepo,
	})

	ctx := context.Background()
	req := request.FulfillmentCallbackRequest{
		ReferenceID:    "ref-123",
		PartnerOrderID: "p-order-123",
		Status:         "SUCCESS",
		Voucher: &request.VoucherData{
			Code: "VOUCHER123",
		},
	}
	idempotencyKey := "idem-456"

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CheckIdempotency", ctx, idempotencyKey).Return(false, "", nil).Once()
		
		order := &aggregate.Order{
			ID:          uuid.New(),
			ReferenceID: "ref-123",
			UserID:      "user-123",
			OrderID:     "order-123",
		}
		mockRepo.On("GetOrderByReferenceID", ctx, "ref-123").Return(order, nil).Once()
		mockRepo.On("CreateFulfillment", ctx, mock.AnythingOfType("*aggregate.Fulfillment")).Return(nil).Once()
		mockRepo.On("UpdateOrder", ctx, mock.AnythingOfType("*aggregate.Order")).Return(nil).Once()
		mockRepo.On("SaveIdempotency", ctx, idempotencyKey, "/sf-backend/fulfillment/callback", 200, "{}").Return(nil).Once()
		
		// For TriggerNotification
		mockRepo.On("GetOrderByReferenceID", ctx, "ref-123").Return(order, nil).Once()
		mockRepo.On("CreateNotification", ctx, mock.AnythingOfType("*aggregate.Notification")).Return(nil).Once()

		err := uc.ProcessFulfillmentCallback(ctx, req, "req-123", idempotencyKey, "sig", "ts")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
