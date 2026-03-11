package wire

import "github.com/leminhthai/train-ticket/booking-service/internal/service"

type App struct {
	BookingService service.BookingService
	PaymentService service.PaymentService
}
