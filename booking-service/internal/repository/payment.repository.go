package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	db "github.com/leminhthai/train-ticket/booking-service/db/generated"
	"github.com/leminhthai/train-ticket/booking-service/internal/model"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, bookingID string, amount float64, method string) error
	CreatePaymentTx(ctx context.Context, params CreatePaymentTxParams) (*CreatePaymentTxResponse, error)
	GetPaymentByBookingID(ctx context.Context, bookingID string) (*db.Payment, error)
	UpdatePaymentStatus(ctx context.Context, bookingID string, status int8) error
}

type paymentRepository struct {
	db *sql.DB
	q  *db.Queries
}

func NewPaymentRepository(db *sql.DB, q *db.Queries) PaymentRepository {
	return &paymentRepository{db: db, q: q}
}

type BookingSeatResult struct {
	SeatID string
	TripID string
}

type CreatePaymentTxParams struct {
	BookingID string
	UserID    string
	Method    string
}

type PaymentResult struct {
	ID        string
	BookingID string
	Amount    float64
	Method    string
	Status    int8
	PaidAt    sql.NullTime
}

type CreatePaymentTxResponse struct {
	Payment      *PaymentResult
	BookingSeats []BookingSeatResult
}

func (pr *paymentRepository) CreatePayment(ctx context.Context, bookingID string, amount float64, method string) error {
	return pr.q.CreatePayment(ctx, db.CreatePaymentParams{
		ID:        uuid.New().String(),
		BookingID: bookingID,
		Amount:    fmt.Sprintf("%.2f", amount),
		Method:    method,
		Status:    int8(model.PaymentStatusPending),
	})
}

func (pr *paymentRepository) GetPaymentByBookingID(ctx context.Context, bookingID string) (*db.Payment, error) {
	payment, err := pr.q.GetPaymentByBookingID(ctx, bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (pr *paymentRepository) UpdatePaymentStatus(ctx context.Context, bookingID string, status int8) error {
	return pr.q.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		Status:    status,
		BookingID: bookingID,
	})
}

func (r *paymentRepository) CreatePaymentTx(ctx context.Context, params CreatePaymentTxParams) (*CreatePaymentTxResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// SELECT FOR UPDATE — lock booking row
	// Chặn double payment race condition
	booking, err := qtx.GetBookingByIDForUpdate(ctx, params.BookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, err
	}

	// Validate ownership
	if booking.UserID != params.UserID {
		return nil, fmt.Errorf("forbidden")
	}

	// Validate status
	switch int(booking.Status) {
	case model.BookingStatusConfirmed:
		return nil, fmt.Errorf("booking already paid")
	case model.BookingStatusCancelled:
		return nil, fmt.Errorf("booking cancelled")
	case model.BookingStatusExpired:
		return nil, fmt.Errorf("booking expired")
	}

	// Validate expiry trong transaction — tránh race condition
	if time.Now().After(booking.ExpiresAt) {
		_ = qtx.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
			Status: int8(model.BookingStatusExpired),
			ID:     params.BookingID,
		})
		_ = tx.Commit()
		return nil, fmt.Errorf("booking expired")
	}

	// Parse total price — handle error
	totalPrice, err := strconv.ParseFloat(booking.TotalPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid booking price")
	}

	// INSERT payment với status = SUCCESS
	paymentID := uuid.New().String()
	if err := qtx.CreatePayment(ctx, db.CreatePaymentParams{
		ID:        paymentID,
		BookingID: params.BookingID,
		Amount:    fmt.Sprintf("%.2f", totalPrice),
		Method:    params.Method,
		Status:    int8(model.PaymentStatusSuccess),
	}); err != nil {
		return nil, err
	}

	// UPDATE booking → CONFIRMED
	if err := qtx.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
		Status: int8(model.BookingStatusConfirmed),
		ID:     params.BookingID,
	}); err != nil {
		return nil, err
	}

	// Lấy booking seats trong transaction
	bookingSeats, err := qtx.GetBookingSeatsByBookingID(ctx, params.BookingID)
	if err != nil {
		return nil, err
	}

	// Lấy payment trong transaction — không cần query lại sau commit
	payment, err := qtx.GetPaymentByBookingID(ctx, params.BookingID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Map sang result types
	amount, _ := strconv.ParseFloat(payment.Amount, 64)

	bookingSeatResults := make([]BookingSeatResult, 0, len(bookingSeats))
	for _, bs := range bookingSeats {
		bookingSeatResults = append(bookingSeatResults, BookingSeatResult{
			SeatID: bs.SeatID,
			TripID: bs.TripID,
		})
	}

	return &CreatePaymentTxResponse{
		Payment: &PaymentResult{
			ID:        payment.ID,
			BookingID: payment.BookingID,
			Amount:    amount,
			Method:    payment.Method,
			Status:    payment.Status,
			PaidAt:    payment.PaidAt,
		},
		BookingSeats: bookingSeatResults,
	}, nil
}
