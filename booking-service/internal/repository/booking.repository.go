package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	db "github.com/leminhthai/train-ticket/booking-service/db/generated"
	"github.com/leminhthai/train-ticket/booking-service/internal/model"
)

type CreateBookingParams struct {
	ID         string
	UserID     string
	TripID     string
	TotalPrice float64
	ExpiresAt  time.Time
}

type CreateBookingTxParams struct {
	BookingID  string
	UserID     string
	TripID     string
	TotalPrice float64
	ExpiresAt  time.Time
	Seats      []model.SeatInfo
}

type CreateBookingSeatParams struct {
	BookingID  string
	SeatID     string
	TripID     string
	SeatNumber string
	Class      string
	Price      float64
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, params CreateBookingParams) error
	CreateBookingTx(ctx context.Context, params CreateBookingTxParams) error
	GetBookingByID(ctx context.Context, bookingID string) (*db.Booking, error)
	GetBookingsByUserID(ctx context.Context, userID string) ([]db.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID string, status int8) error
	GetExpiredBookings(ctx context.Context) ([]db.Booking, error)
	CreateBookingSeat(ctx context.Context, params CreateBookingSeatParams) error
	GetBookingSeatsByBookingID(ctx context.Context, bookingID string) ([]db.BookingSeat, error)
	GetBookedSeatsByTripID(ctx context.Context, tripId string) ([]db.GetBookedSeatsByTripIDRow, error)
	IsSeatBooked(ctx context.Context, tripID, seatID string) (bool, error)
	AreSeatsBooked(ctx context.Context, tripID string, seatIDs []string) ([]string, error)
}

type bookingRepository struct {
	db *sql.DB
	q  *db.Queries
}

func NewBookingRepository(db *sql.DB, q *db.Queries) BookingRepository {
	return &bookingRepository{
		db: db,
		q:  q,
	}
}

func (br *bookingRepository) CreateBooking(ctx context.Context, params CreateBookingParams) error {

	return br.q.CreateBooking(ctx, db.CreateBookingParams{
		ID:         params.ID,
		UserID:     params.UserID,
		TripID:     params.TripID,
		TotalPrice: fmt.Sprintf("%.2f", params.TotalPrice),
		Status:     int8(model.BookingStatusPending),
		ExpiresAt:  params.ExpiresAt,
	})
}

func (br *bookingRepository) GetBookingByID(ctx context.Context, bookingID string) (*db.Booking, error) {
	booking, err := br.q.GetBookingByID(ctx, bookingID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, err
	}

	return &booking, nil
}

func (br *bookingRepository) GetBookingsByUserID(ctx context.Context, userID string) ([]db.Booking, error) {
	return br.q.GetBookingsByUserID(ctx, userID)
}

func (br *bookingRepository) UpdateBookingStatus(ctx context.Context, bookingID string, status int8) error {
	return br.q.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
		ID:     bookingID,
		Status: status,
	})
}

func (br *bookingRepository) GetExpiredBookings(ctx context.Context) ([]db.Booking, error) {
	return br.q.GetExpiredBookings(ctx)
}

func (br *bookingRepository) CreateBookingSeat(ctx context.Context, params CreateBookingSeatParams) error {
	return br.q.CreateBookingSeat(ctx, db.CreateBookingSeatParams{
		ID:         uuid.New().String(),
		BookingID:  params.BookingID,
		SeatID:     params.SeatID,
		TripID:     params.TripID,
		SeatNumber: params.SeatNumber,
		Class:      params.Class,
		Price:      fmt.Sprintf("%.2f", params.Price),
	})
}

func (br *bookingRepository) GetBookingSeatsByBookingID(ctx context.Context, bookingID string) ([]db.BookingSeat, error) {
	return br.q.GetBookingSeatsByBookingID(ctx, bookingID)
}

func (br *bookingRepository) GetBookedSeatsByTripID(ctx context.Context, tripID string) ([]db.GetBookedSeatsByTripIDRow, error) {
	return br.q.GetBookedSeatsByTripID(ctx, tripID)
}

func (br *bookingRepository) CreateBookingTx(ctx context.Context, params CreateBookingTxParams) error {
	tx, err := br.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := br.q.WithTx(tx)

	// SELECT FOR UPDATE — lock rows trong transaction
	bookedSeats, err := qtx.GetBookedSeatsByTripIDForUpdate(ctx, params.TripID)
	if err != nil {
		return err
	}

	bookedMap := make(map[string]bool)
	for _, seatID := range bookedSeats {
		bookedMap[seatID] = true
	}
	for _, seat := range params.Seats {
		if bookedMap[seat.SeatID] {
			return fmt.Errorf("seat %s is already booked", seat.SeatID)
		}
	}

	// Insert booking
	if err := qtx.CreateBooking(ctx, db.CreateBookingParams{
		ID:         params.BookingID,
		UserID:     params.UserID,
		TripID:     params.TripID,
		TotalPrice: fmt.Sprintf("%.2f", params.TotalPrice),
		Status:     int8(model.BookingStatusPending),
		ExpiresAt:  params.ExpiresAt,
	}); err != nil {
		return err
	}

	// Insert booking_seats
	for _, seat := range params.Seats {
		if err := qtx.CreateBookingSeat(ctx, db.CreateBookingSeatParams{
			ID:         uuid.New().String(),
			BookingID:  params.BookingID,
			SeatID:     seat.SeatID,
			TripID:     params.TripID,
			SeatNumber: seat.SeatNumber,
			Class:      seat.Class,
			Price:      fmt.Sprintf("%.2f", seat.Price),
		}); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *bookingRepository) IsSeatBooked(ctx context.Context, tripID, seatID string) (bool, error) {
	rows, err := r.q.GetBookedSeatsByTripID(ctx, tripID)
	if err != nil {
		return false, err
	}
	for _, row := range rows {
		if row.SeatID == seatID {
			return true, nil
		}
	}
	return false, nil
}

func (r *bookingRepository) AreSeatsBooked(ctx context.Context, tripID string, seatIDs []string) ([]string, error) {
	if len(seatIDs) == 0 {
		return nil, nil
	}

	// Build query động với IN clause
	query := `
		SELECT bs.seat_id FROM booking_seats bs
		JOIN bookings b ON bs.booking_id = b.id
		WHERE bs.trip_id = ?
		AND bs.seat_id IN (?` + strings.Repeat(",?", len(seatIDs)-1) + `)
		AND b.status IN (1, 2)
	`

	args := make([]interface{}, 0, len(seatIDs)+1)
	args = append(args, tripID)
	for _, id := range seatIDs {
		args = append(args, id)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	booked := make([]string, 0)
	for rows.Next() {
		var seatID string
		if err := rows.Scan(&seatID); err != nil {
			return nil, err
		}
		booked = append(booked, seatID)
	}

	return booked, nil
}
