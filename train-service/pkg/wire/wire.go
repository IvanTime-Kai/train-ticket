//go:build wireinject
// +build wireinject

package wire

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/leminhthai/train-ticket/train-service/internal/repository"
	"github.com/leminhthai/train-ticket/train-service/internal/service"
)

// StationSet nhóm provider cho station (repo -> service). *db.Queries do ProvideQueries cung cấp một lần ở InitializeApp.
var StationSet = wire.NewSet(
	repository.NewStationRepository,
	service.NewStationService,
)

// TrainSet nhóm provider cho train (repo -> service).
var TrainSet = wire.NewSet(
	repository.NewTrainRepository,
	service.NewTrainService,
)

// TripSet nhóm provider cho trip (repo -> service).
var TripSet = wire.NewSet(
	repository.NewTripRepository,
	service.NewTripService,
)

// InitializeApp tạo toàn bộ dependency. ProvideQueries chỉ gọi một lần, tránh lỗi "multiple bindings".
func InitializeApp(sqlDB *sql.DB) *App {
	wire.Build(
		ProvideQueries,
		StationSet,
		TrainSet,
		TripSet,
		wire.Struct(new(App), "*"),
	)
	return nil
}
