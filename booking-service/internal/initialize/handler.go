package initialize

import (
    "github.com/leminhthai/train-ticket/booking-service/internal/handler"
    "github.com/leminhthai/train-ticket/booking-service/pkg/wire"
)

func InitHandlers(app *wire.App) {
    handler.InitBooking(app.BookingService)
    handler.InitPayment(app.PaymentService)
}