package service

import (
	"context"
	"fmt"

	db "github.com/leminhthai/train-ticket/train-service/db/generated"
	"github.com/leminhthai/train-ticket/train-service/internal/model"
	"github.com/leminhthai/train-ticket/train-service/internal/repository"
)

type TrainService interface {
	Create(ctx context.Context, req *model.CreateTrainRequest) (db.Train, error)
	GetByID(ctx context.Context, id string) (db.Train, error)
	List(ctx context.Context) ([]db.Train, error)
	DeActive(ctx context.Context, id string) error
	AddSeat(ctx context.Context, trainID string, req *model.CreateSeatRequest) (db.Seat, error)
	ListSeat(ctx context.Context, trainID string) ([]db.Seat, error)
}

type trainService struct {
	repo repository.TrainRepository
}

func NewTrainService(repo repository.TrainRepository) TrainService {
	return &trainService{repo: repo}
}

func (ts trainService) Create(ctx context.Context, req *model.CreateTrainRequest) (db.Train, error) {
	return ts.repo.Create(ctx, req)
}

func (ts trainService) GetByID(ctx context.Context, id string) (db.Train, error) {
	train, err := ts.repo.GetByID(ctx, id)
	if err != nil {
		return db.Train{}, fmt.Errorf("train not found")
	}

	return train, nil
}

func (ts *trainService) List(ctx context.Context) ([]db.Train, error) {
	return ts.repo.List(ctx)
}

func (ts *trainService) DeActive(ctx context.Context, id string) error {
	_, err := ts.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("train not found")
	}

	return ts.repo.UpdateStatus(ctx, id, model.TrainStatusInactive)
}

func (ts *trainService) AddSeat(ctx context.Context, trainID string, req *model.CreateSeatRequest) (db.Seat, error) {
	_, err := ts.repo.GetByID(ctx, trainID)
	if err != nil {
		return db.Seat{}, fmt.Errorf("train not found")
	}
	return ts.repo.CreateSeat(ctx, trainID, req)
}

func (ts *trainService) ListSeat(ctx context.Context, trainID string) ([]db.Seat, error) {
	return ts.repo.ListSeats(ctx, trainID)
}