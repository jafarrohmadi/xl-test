package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xlsmart-api/sf-payment/handler/exceptions"
	"github.com/xlsmart-api/sf-payment/handler/validate"
	"github.com/xlsmart-api/sf-payment/model/request"
	"github.com/xlsmart-api/sf-payment/model/response"
)

// PostSfPaymentPaymentsRequest handles payment request from SF Backend
func (s *Server) PostSfPaymentPaymentsRequest(ctx echo.Context) error {
	var (
		context        = ctx.Request().Context()
		requestData    request.PaymentRequest
		requestID      = ctx.Request().Header.Get("X-Request-Id")
		idempotencyKey = ctx.Request().Header.Get("Idempotency-Key")
	)

	if err := ctx.Bind(&requestData); err != nil {
		httpResponse := response.BuildErrorResponse(response.InvalidRequest, requestID, nil)
		return ctx.JSON(http.StatusBadRequest, httpResponse)
	}

	if err := validate.ValidateStruct(requestData); err != nil {
		httpResponse := response.BuildErrorResponse(response.InvalidRequest, requestID, nil)
		return ctx.JSON(http.StatusBadRequest, httpResponse)
	}

	data, err := s.UseCase.ProcessPayment(context, requestData, requestID, idempotencyKey)
	if err != nil {
		httpResponse := exceptions.HandleError(context, err, requestID)
		return ctx.JSON(httpResponse.HTTPStatusCode, httpResponse)
	}

	httpResponse := response.BuildSuccessResponse(response.PaymentSuccess, requestID, data)
	return ctx.JSON(http.StatusOK, httpResponse)
}

// PostSfPaymentPaymentsWebhook handles payment gateway webhook
func (s *Server) PostSfPaymentPaymentsWebhook(ctx echo.Context) error {
	var (
		context        = ctx.Request().Context()
		requestData    request.PaymentWebhookRequest
		requestID      = ctx.Request().Header.Get("X-Request-Id")
		idempotencyKey = ctx.Request().Header.Get("Idempotency-Key")
		signature      = ctx.Request().Header.Get("X-Signature")
		timestamp      = ctx.Request().Header.Get("X-Signature-Timestamp")
	)

	if err := ctx.Bind(&requestData); err != nil {
		httpResponse := response.BuildErrorResponse(response.InvalidRequest, requestID, nil)
		return ctx.JSON(http.StatusBadRequest, httpResponse)
	}

	if err := validate.ValidateStruct(requestData); err != nil {
		httpResponse := response.BuildErrorResponse(response.InvalidRequest, requestID, nil)
		return ctx.JSON(http.StatusBadRequest, httpResponse)
	}

	err := s.UseCase.ProcessWebhook(context, requestData, requestID, idempotencyKey, signature, timestamp)
	if err != nil {
		httpResponse := exceptions.HandleError(context, err, requestID)
		return ctx.JSON(httpResponse.HTTPStatusCode, httpResponse)
	}

	httpResponse := response.BuildSuccessResponse(response.WebhookAccepted, requestID, nil)
	return ctx.JSON(http.StatusOK, httpResponse)
}
