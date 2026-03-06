package service

import (
	"context"
	"fmt"

	db "github.com/leminhthai/train-ticket/train-service/db/generated"
	"github.com/leminhthai/train-ticket/train-service/internal/model"
	"github.com/leminhthai/train-ticket/train-service/internal/repository"
)

type StationService interface {
	Create(ctx context.Context, req *model.CreateStationRequest) (db.Station, error)
	GetById(ctx context.Context, id string) (db.Station, error)
	GetByCode(ctx context.Context, code string) (db.Station, error)
	List(ctx context.Context) ([]db.Station, error)
}

type stationService struct {
	repo repository.StationRepository
}

func NewStationService(repo repository.StationRepository) StationService {
	return &stationService{repo: repo}
}

func (ss *stationService) Create(ctx context.Context, req *model.CreateStationRequest) (db.Station, error) {
	// check station existing
	existing, err := ss.repo.GetByCode(ctx, req.Code)

	if err == nil || existing.ID != "" {
		return db.Station{}, fmt.Errorf("station code %s already exist", req.Code)
	}

	return ss.repo.Create(ctx, req)
}

func (ss *stationService) GetById(ctx context.Context, id string) (db.Station, error) {
	station, err := ss.repo.GetById(ctx, id)
	if err != nil {
		return db.Station{}, fmt.Errorf("station not found")
	}
	return station, nil
}

func (ss *stationService) GetByCode(ctx context.Context, code string) (db.Station, error) {
	station, err := ss.repo.GetByCode(ctx, code)
	if err != nil {
		return db.Station{}, fmt.Errorf("station not found")
	}
	return station, nil
}

func (ss *stationService) List(ctx context.Context) ([]db.Station, error) {
	return ss.repo.List(ctx)
}
