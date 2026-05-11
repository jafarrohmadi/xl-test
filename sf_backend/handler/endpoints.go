package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xlsmart-api/sf-backend/handler/exceptions"
	"github.com/xlsmart-api/sf-backend/handler/validate"
	"github.com/xlsmart-api/sf-backend/model/request"
	"github.com/xlsmart-api/sf-backend/model/response"
)

// PostOrders handles order submission - POST /orders
func (s *Server) PostOrders(ctx echo.Context) error {
	var (
		context        = ctx.Request().Context()
		requestData    request.SubmitOrderRequest
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

	data, err := s.UseCase.SubmitOrder(context, requestData, requestID, idempotencyKey)
	if err != nil {
		httpResponse := exceptions.HandleError(context, err, requestID)
		return ctx.JSON(httpResponse.HTTPStatusCode, httpResponse)
	}

	httpResponse := response.BuildSuccessResponse(response.OrderCreated, requestID, data)
	return ctx.JSON(http.StatusCreated, httpResponse)
}

// PostOrdersFulfillmentCallback handles fulfillment callback from Partner - POST /orders/fulfillment/callback
func (s *Server) PostOrdersFulfillmentCallback(ctx echo.Context) error {
	var (
		context        = ctx.Request().Context()
		requestData    request.FulfillmentCallbackRequest
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

	err := s.UseCase.ProcessFulfillmentCallback(context, requestData, requestID, idempotencyKey, signature, timestamp)
	if err != nil {
		httpResponse := exceptions.HandleError(context, err, requestID)
		return ctx.JSON(httpResponse.HTTPStatusCode, httpResponse)
	}

	httpResponse := response.BuildSuccessResponse(response.FulfillmentCallbackAccepted, requestID, nil)
	return ctx.JSON(http.StatusOK, httpResponse)
}

// PostInternalNotifications handles notification trigger - POST /internal/notifications
func (s *Server) PostInternalNotifications(ctx echo.Context) error {
	var (
		context     = ctx.Request().Context()
		requestData request.NotificationEvent
		requestID   = ctx.Request().Header.Get("X-Request-Id")
	)

	if err := ctx.Bind(&requestData); err != nil {
		httpResponse := response.BuildErrorResponse(response.InvalidRequest, requestID, nil)
		return ctx.JSON(http.StatusBadRequest, httpResponse)
	}

	if err := validate.ValidateStruct(requestData); err != nil {
		httpResponse := response.BuildErrorResponse(response.InvalidRequest, requestID, nil)
		return ctx.JSON(http.StatusBadRequest, httpResponse)
	}

	err := s.UseCase.TriggerNotification(context, requestData, requestID)
	if err != nil {
		httpResponse := exceptions.HandleError(context, err, requestID)
		return ctx.JSON(httpResponse.HTTPStatusCode, httpResponse)
	}

	httpResponse := response.BuildSuccessResponse(response.NotificationQueued, requestID, nil)
	return ctx.JSON(http.StatusAccepted, httpResponse)
}
