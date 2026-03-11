package router

type RouterGroup struct {
	BookingRouter BookingRouter
	PaymentRouter PaymentRouter
}

var RouterGroupApp = new(RouterGroup)
