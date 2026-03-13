package wire

import (
	"github.com/leminhthai/train-ticket/booking-service/internal/repository"
	"github.com/leminhthai/train-ticket/booking-service/internal/service"
)

type App struct {
	BookingService    service.BookingService
	PaymentService    service.PaymentService
	BookingRepository repository.BookingRepository
}
