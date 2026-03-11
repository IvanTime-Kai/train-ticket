//go:build wireinject
// +build wireinject
package wire

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/leminhthai/train-ticket/booking-service/internal/repository"
	"github.com/leminhthai/train-ticket/booking-service/internal/service"
)

var BookingSet = wire.NewSet(
	repository.NewBookingRepository,
	service.NewBookingService,
)

var PaymentSet = wire.NewSet(
	repository.NewPaymentRepository,
	service.NewPaymentService,
)

func InitializeApp(sqlDB *sql.DB) *App {
	wire.Build(
		ProvideQueries,
		BookingSet,
		PaymentSet,
		wire.Struct(new(App), "*"),
	)

	return nil
}
