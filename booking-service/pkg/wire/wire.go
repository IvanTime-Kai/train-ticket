//go:build wireinject
// +build wireinject
package wire

import (
	"database/sql"

	"github.com/google/wire"
	grpcClient "github.com/leminhthai/train-ticket/booking-service/internal/grpc"
	"github.com/leminhthai/train-ticket/booking-service/internal/repository"
	"github.com/leminhthai/train-ticket/booking-service/internal/service"
	"github.com/leminhthai/train-ticket/booking-service/pkg/setting"
)

var BookingSet = wire.NewSet(
	repository.NewBookingRepository,
	service.NewBookingService,
	ProvideTrainClient,
)

var PaymentSet = wire.NewSet(
	repository.NewPaymentRepository,
	service.NewPaymentService,
)

func ProvideTrainClient(cfg setting.GRPCSetting) (*grpcClient.TrainClient, error) {
	return grpcClient.NewTrainClient(cfg.TrainServiceHost, cfg.TrainServicePort)
}

func InitializeApp(sqlDB *sql.DB, grpcCfg setting.GRPCSetting) (*App, error) {
	wire.Build(
		ProvideQueries,
		BookingSet,
		PaymentSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
