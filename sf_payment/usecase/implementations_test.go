package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xlsmart-api/sf-payment/model/aggregate"
	"github.com/xlsmart-api/sf-payment/model/request"
	repo_mocks "github.com/xlsmart-api/sf-payment/repository/mocks"
)

func TestProcessPayment(t *testing.T) {
	mockRepo := new(repo_mocks.MockRepository)
	uc := NewUseCase(NewUseCaseOptions{
		Repository: mockRepo,
	})

	ctx := context.Background()
	req := request.PaymentRequest{
		ReferenceID: "ref-123",
		Price:       100000,
	}
	requestID := "req-123"
	idempotencyKey := "idem-123"

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CheckIdempotency", ctx, idempotencyKey).Return(false, "", nil).Once()
		mockRepo.On("CreatePaymentTransaction", ctx, mock.AnythingOfType("*aggregate.PaymentTransaction")).Return(nil).Once()
		mockRepo.On("UpdatePaymentTransaction", ctx, mock.AnythingOfType("*aggregate.PaymentTransaction")).Return(nil).Once()
		mockRepo.On("SaveIdempotency", ctx, idempotencyKey, "/sf-payment/payments/request", 200, mock.Anything).Return(nil).Once()

		resp, err := uc.ProcessPayment(ctx, req, requestID, idempotencyKey)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "ref-123", resp["referenceId"])
		assert.Equal(t, "SUCCESS", resp["paymentStatus"])
		mockRepo.AssertExpectations(t)
	})
}

func TestProcessWebhook(t *testing.T) {
	mockRepo := new(repo_mocks.MockRepository)
	uc := NewUseCase(NewUseCaseOptions{
		Repository: mockRepo,
	})

	ctx := context.Background()
	req := request.PaymentWebhookRequest{
		ReferenceID:   "ref-123",
		TransactionID: "trx-456",
		Status:        "SUCCESS",
		Amount:        100000,
	}
	idempotencyKey := "idem-456"
	timestamp := time.Now().Format(time.RFC3339)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CheckIdempotency", ctx, idempotencyKey).Return(false, "", nil).Once()
		
		payment := &aggregate.PaymentTransaction{
			ID:          uuid.New(),
			ReferenceID: "ref-123",
			Amount:      100000,
		}
		mockRepo.On("GetPaymentByReferenceID", ctx, "ref-123").Return(payment, nil).Once()
		mockRepo.On("CreatePaymentWebhook", ctx, mock.AnythingOfType("*aggregate.PaymentWebhook")).Return(nil).Once()
		mockRepo.On("UpdatePaymentTransaction", ctx, mock.AnythingOfType("*aggregate.PaymentTransaction")).Return(nil).Once()
		mockRepo.On("SaveIdempotency", ctx, idempotencyKey, "/sf-payment/payments/webhook", 200, "{}").Return(nil).Once()

		err := uc.ProcessWebhook(ctx, req, "req-123", idempotencyKey, "sig", timestamp)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
