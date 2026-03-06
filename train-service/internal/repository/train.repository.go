package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	db "github.com/leminhthai/train-ticket/train-service/db/generated"
	"github.com/leminhthai/train-ticket/train-service/internal/model"
)

type TrainRepository interface {
	Create(ctx context.Context, req *model.CreateTrainRequest) (db.Train, error)
	GetByID(ctx context.Context, id string) (db.Train, error)
	List(ctx context.Context) ([]db.Train, error)
	UpdateStatus(ctx context.Context, id string, status int) error
	CreateSeat(ctx context.Context, trainID string, req *model.CreateSeatRequest) (db.Seat, error)
	ListSeats(ctx context.Context, trainID string) ([]db.Seat, error)
	GetSeatByID(ctx context.Context, id string) (db.Seat, error)
}

type trainRepository struct {
	queries *db.Queries
}

func NewTrainRepository(queries *db.Queries) TrainRepository {
	return &trainRepository{queries: queries}
}

func (r *trainRepository) Create(ctx context.Context, req *model.CreateTrainRequest) (db.Train, error) {
	id := uuid.New().String()

	err := r.queries.CreateTrain(ctx, db.CreateTrainParams{
		ID:         id,
		Name:       req.Name,
		TotalSeats: int32(req.TotalSeats),
		Status:     model.TrainStatusActive,
	})
	if err != nil {
		return db.Train{}, err
	}

	return r.queries.GetTrainByID(ctx, id)
}

func (r *trainRepository) GetByID(ctx context.Context, id string) (db.Train, error) {
	return r.queries.GetTrainByID(ctx, id)
}

func (r *trainRepository) List(ctx context.Context) ([]db.Train, error) {
	return r.queries.ListTrains(ctx)
}

func (r *trainRepository) UpdateStatus(ctx context.Context, id string, status int) error {
	return r.queries.UpdateTrainStatus(ctx, db.UpdateTrainStatusParams{
		ID:     id,
		Status: int8(status),
	})
}

func (r *trainRepository) CreateSeat(ctx context.Context, trainID string, req *model.CreateSeatRequest) (db.Seat, error) {
	id := uuid.New().String()

	err := r.queries.CreateSeat(ctx, db.CreateSeatParams{
		ID:         id,
		TrainID:    trainID,
		SeatNumber: req.SeatNumber,
		Class:      req.Class,
		Price:      fmt.Sprintf("%.2f", req.Price),
	})
	if err != nil {
		return db.Seat{}, err
	}

	return r.queries.GetSeatByID(ctx, id)
}

func (r *trainRepository) ListSeats(ctx context.Context, trainID string) ([]db.Seat, error) {
	return r.queries.ListSeatsByTrain(ctx, trainID)
}

func (r *trainRepository) GetSeatByID(ctx context.Context, id string) (db.Seat, error) {
	return r.queries.GetSeatByID(ctx, id)
}