package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	db "github.com/leminhthai/train-ticket/train-service/db/generated"
	"github.com/leminhthai/train-ticket/train-service/internal/model"
	"github.com/leminhthai/train-ticket/train-service/internal/utils"
)

type TripRepository interface {
	Create(ctx context.Context, req *model.CreateTripRequest) (db.Trip, error)
	GetByID(ctx context.Context, id string) (db.Trip, error)
	Search(ctx context.Context, req *model.SearchTripRequest) ([]db.Trip, error)
	UpdateStatus(ctx context.Context, id string, status int) error
	CreateRoute(ctx context.Context, req *model.CreateRouteRequest) (db.Route, error)
	GetRouteByID(ctx context.Context, id string) (db.Route, error)
}

type tripRepository struct {
	queries *db.Queries
}

func NewTripRepository(queries *db.Queries) TripRepository {
	return &tripRepository{queries: queries}
}

func (r *tripRepository) Create(ctx context.Context, req *model.CreateTripRequest) (db.Trip, error) {
	id := uuid.New().String()

	departureTime, err := time.Parse(utils.DateTimeFormat, req.DepartureTime)
	if err != nil {
		return db.Trip{}, fmt.Errorf("invalid departure_time format")
	}

	arrivalTime, err := time.Parse(utils.DateTimeFormat, req.ArrivalTime)
	if err != nil {
		return db.Trip{}, fmt.Errorf("invalid arrival_time format")
	}

	err = r.queries.CreateTrip(ctx, db.CreateTripParams{
		ID:            id,
		TrainID:       req.TrainID,
		RouteID:       req.RouteID,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		Status:        1,
	})
	if err != nil {
		return db.Trip{}, err
	}

	return r.queries.GetTripByID(ctx, id)
}

func (r *tripRepository) GetByID(ctx context.Context, id string) (db.Trip, error) {
	return r.queries.GetTripByID(ctx, id)
}

func (r *tripRepository) Search(ctx context.Context, req *model.SearchTripRequest) ([]db.Trip, error) {
	// 1. Get station by code
	origin, err := r.queries.GetStationByCode(ctx, req.From)
	if err != nil {
		return nil, fmt.Errorf("origin station not found: %s", req.From)
	}

	destination, err := r.queries.GetStationByCode(ctx, req.To)
	if err != nil {
		return nil, fmt.Errorf("destination station not found: %s", req.To)
	}

	// Parse date
	date, err := time.Parse(utils.DateFormat, req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}

	// 3. cal range: 00:00:00 → 00:00:00 (tomorrow)
	startOfDay := date
	endOfDay := date.Add(24 * time.Hour)

	return r.queries.SearchTrips(ctx, db.SearchTripsParams{
		OriginStationID:      origin.ID,
		DestinationStationID: destination.ID,
		DepartureTime:        startOfDay, // start: 00:00:00 current date
		DepartureTime_2:      endOfDay,   // end: 00:00:00 tomorrow
	})
}

func (r *tripRepository) UpdateStatus(ctx context.Context, id string, status int) error {
	return r.queries.UpdateTripStatus(ctx, db.UpdateTripStatusParams{
		ID:     id,
		Status: int8(status),
	})
}

func (r *tripRepository) CreateRoute(ctx context.Context, req *model.CreateRouteRequest) (db.Route, error) {
	id := uuid.New().String()
	err := r.queries.CreateRoute(ctx, db.CreateRouteParams{
		ID:                   id,
		OriginStationID:      req.OriginStationID,
		DestinationStationID: req.DestinationStationID,
		DistanceKm:           sql.NullInt32{Int32: int32(req.DistanceKm), Valid: req.DistanceKm > 0},
	})
	if err != nil {
		return db.Route{}, err
	}
	return r.queries.GetRouteByID(ctx, id)
}

func (r *tripRepository) GetRouteByID(ctx context.Context, id string) (db.Route, error) {
	return r.queries.GetRouteByID(ctx, id)
}
