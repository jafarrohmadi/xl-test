package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xlsmart-api/partner/model/aggregate"
	"github.com/xlsmart-api/partner/model/request"
	repo_mocks "github.com/xlsmart-api/partner/repository/mocks"
)

func TestSubmitOrder(t *testing.T) {
	mockRepo := new(repo_mocks.MockRepository)
	uc := NewUseCase(NewUseCaseOptions{
		Repository: mockRepo,
	})

	ctx := context.Background()
	req := request.PartnerSubmitRequest{
		OrderID:     "order-123",
		ReferenceID: "ref-123",
		TotalPrice:  100000,
	}
	requestID := "req-123"
	idempotencyKey := "idem-123"

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CheckIdempotency", ctx, idempotencyKey).Return(false, "", nil).Once()
		mockRepo.On("CreatePartnerOrder", ctx, mock.AnythingOfType("*aggregate.Order")).Return(nil).Once()
		mockRepo.On("SaveIdempotency", ctx, idempotencyKey, "/partner/orders/submit", 200, mock.Anything).Return(nil).Once()

		resp, err := uc.SubmitOrder(ctx, req, requestID, idempotencyKey)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Contains(t, resp["partnerOrderId"], "P-")
		assert.Equal(t, "ACCEPTED", resp["status"])
		mockRepo.AssertExpectations(t)
	})
}

func TestProcessFulfillment(t *testing.T) {
	mockRepo := new(repo_mocks.MockRepository)
	uc := NewUseCase(NewUseCaseOptions{
		Repository: mockRepo,
	})

	ctx := context.Background()
	req := request.FulfillmentRequest{
		ReferenceID:    "ref-123",
		PartnerOrderID: "P-12345",
	}
	idempotencyKey := "idem-456"

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CheckIdempotency", ctx, idempotencyKey).Return(false, "", nil).Once()
		
		order := &aggregate.Order{
			ID:             uuid.New(),
			ReferenceID:    "ref-123",
			PartnerOrderID: "P-12345",
		}
		mockRepo.On("GetPartnerOrderByReferenceID", ctx, "ref-123").Return(order, nil).Once()
		mockRepo.On("CreateFulfillment", ctx, mock.AnythingOfType("*aggregate.Fulfillment")).Return(nil).Once()
		mockRepo.On("SaveIdempotency", ctx, idempotencyKey, "/partner/orders/fulfillment", 200, "{}").Return(nil).Once()

		err := uc.ProcessFulfillment(ctx, req, "req-123", idempotencyKey)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
