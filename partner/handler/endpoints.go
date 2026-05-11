package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xlsmart-api/partner/handler/exceptions"
	"github.com/xlsmart-api/partner/handler/validate"
	"github.com/xlsmart-api/partner/model/request"
	"github.com/xlsmart-api/partner/model/response"
)

// PostPartnerOrdersSubmit handles order submission to Partner
func (s *Server) PostPartnerOrdersSubmit(ctx echo.Context) error {
	var (
		context        = ctx.Request().Context()
		requestData    request.PartnerSubmitRequest
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

	httpResponse := response.BuildSuccessResponse(response.PartnerOrderAccepted, requestID, data)
	return ctx.JSON(http.StatusOK, httpResponse)
}

// PostPartnerOrdersFulfillment handles fulfillment request
func (s *Server) PostPartnerOrdersFulfillment(ctx echo.Context) error {
	var (
		context        = ctx.Request().Context()
		requestData    request.FulfillmentRequest
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

	err := s.UseCase.ProcessFulfillment(context, requestData, requestID, idempotencyKey)
	if err != nil {
		httpResponse := exceptions.HandleError(context, err, requestID)
		return ctx.JSON(httpResponse.HTTPStatusCode, httpResponse)
	}

	httpResponse := response.BuildSuccessResponse(response.FulfillmentInProgress, requestID, nil)
	return ctx.JSON(http.StatusOK, httpResponse)
}
