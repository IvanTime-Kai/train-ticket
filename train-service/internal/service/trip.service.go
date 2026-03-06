package service

import (
	"context"
	"fmt"

	db "github.com/leminhthai/train-ticket/train-service/db/generated"
	"github.com/leminhthai/train-ticket/train-service/internal/model"
	"github.com/leminhthai/train-ticket/train-service/internal/repository"
)

type TripService interface {
	Create(ctx context.Context, req *model.CreateTripRequest) (db.Trip, error)
	GetByID(ctx context.Context, id string) (db.Trip, error)
	Search(ctx context.Context, req *model.SearchTripRequest) ([]db.Trip, error)
	Cancel(ctx context.Context, id string) error
	CreateRoute(ctx context.Context, req *model.CreateRouteRequest) (db.Route, error)
	GetRouteByID(ctx context.Context, id string) (db.Route, error)
}

type tripService struct {
	repo repository.TripRepository
}

func NewTripService(repo repository.TripRepository) TripService {
	return &tripService{repo: repo}
}

func (ts *tripService) Create(ctx context.Context, req *model.CreateTripRequest) (db.Trip, error) {
	return ts.repo.Create(ctx, req)
}

func (ts *tripService) GetByID(ctx context.Context, id string) (db.Trip, error) {
	trip, err := ts.repo.GetByID(ctx, id)
	if err != nil {
		return db.Trip{}, fmt.Errorf("trip not found")
	}
	return trip, nil
}

func (ts *tripService) Search(ctx context.Context, req *model.SearchTripRequest) ([]db.Trip, error) {
	return ts.repo.Search(ctx, req)
}

func (ts *tripService) Cancel(ctx context.Context, id string) error {
	_, err := ts.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("trip not found")
	}
	return ts.repo.UpdateStatus(ctx, id, model.TripStatusCancelled)
}

func (ts *tripService) CreateRoute(ctx context.Context, req *model.CreateRouteRequest) (db.Route, error) {
	return ts.repo.CreateRoute(ctx, req)
}

func (ts *tripService) GetRouteByID(ctx context.Context, id string) (db.Route, error) {
	return ts.repo.GetRouteByID(ctx, id)
}