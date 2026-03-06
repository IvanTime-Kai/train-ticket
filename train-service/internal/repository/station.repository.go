package repository

import (
	"context"

	"github.com/google/uuid"
	db "github.com/leminhthai/train-ticket/train-service/db/generated"
	"github.com/leminhthai/train-ticket/train-service/internal/model"
)

type StationRepository interface {
	Create(ctx context.Context, req *model.CreateStationRequest) (db.Station, error)
	GetById(ctx context.Context, id string) (db.Station, error)
	GetByCode(ctx context.Context, code string) (db.Station, error)
	List(cx context.Context) ([]db.Station, error)
}

type stationRepository struct {
	queries *db.Queries
}

func NewStationRepository(queries *db.Queries) StationRepository {
	return &stationRepository{
		queries: queries,
	}
}

func (us *stationRepository) Create(ctx context.Context, req *model.CreateStationRequest) (db.Station, error) {

	id := uuid.New().String()

	err := us.queries.CreateStation(ctx, db.CreateStationParams{
		ID: id,
		Name: req.Name,
		Code: req.Code,
		City: req.City,
	})

	if err != nil {
		return db.Station{}, err
	}

	return us.queries.GetStationByID(ctx, id)
}

func (us *stationRepository) GetById(ctx context.Context, id string) (db.Station, error) {
	return us.queries.GetStationByID(ctx, id)
}

func (us *stationRepository) GetByCode(ctx context.Context, code string) (db.Station, error) {
	return us.queries.GetStationByCode(ctx, code)
}

func (us *stationRepository) List(ctx context.Context) ([]db.Station, error) {
	return us.queries.ListStations(ctx)
}
