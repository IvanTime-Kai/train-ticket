package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leminhthai/train-ticket/booking-service/internal/handler"
)

type PaymentRouter struct{}

func (pr *PaymentRouter) InitPaymentRouter(Router *gin.RouterGroup) {

	// PRIVATE
	paymentPrivate := Router.Group("/payments")
	{
		paymentPrivate.POST("", handler.Payment.CreatePayment)
	}
}
