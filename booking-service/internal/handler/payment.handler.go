package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/booking-service/internal/model"
	"github.com/leminhthai/train-ticket/booking-service/internal/service"
	"github.com/leminhthai/train-ticket/booking-service/pkg/response"
)

var Payment = new(PaymentHandler)

type PaymentHandler struct {
	ps service.PaymentService
}

func InitPayment(ps service.PaymentService) {
	Payment.ps = ps
}

// CreatePayment godoc
// @Summary     Thanh toán booking
// @Description Thanh toán cho booking đang pending
// @Tags        payments
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body model.CreatePaymentRequest true "Create Payment Request"
// @Success     200 {object} response.ResponseData
// @Failure     400 {object} response.ResponseData
// @Failure     401 {object} response.ResponseData
// @Router      /payments [post]
func (bh *PaymentHandler) CreatePayment(c *gin.Context) {
	var req model.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	userID := c.GetString("subjectUUID")

	result, err := bh.ps.CreatePayment(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorResponse(c, response.ErrCodeParamsInvalid, err.Error())
		return
	}

	response.SuccessResponse(c, response.ErrCodeSuccess, result)
}
