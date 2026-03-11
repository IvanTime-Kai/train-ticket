package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/leminhthai/train-ticket/booking-service/internal/model"
	"github.com/leminhthai/train-ticket/booking-service/internal/repository"
	"github.com/leminhthai/train-ticket/booking-service/internal/utils/cache"
	timeFormat "github.com/leminhthai/train-ticket/booking-service/internal/utils/time"
	"go.uber.org/zap"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, userID string, req *model.CreatePaymentRequest) (*model.PaymentResponse, error)
}

type paymentService struct {
	bookingRepo repository.BookingRepository
	paymentRepo repository.PaymentRepository
}

func NewPaymentService(bookingRepo repository.BookingRepository, paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentService{
		bookingRepo: bookingRepo,
		paymentRepo: paymentRepo,
	}
}

func (ps *paymentService) CreatePayment(ctx context.Context, userID string, req *model.CreatePaymentRequest) (*model.PaymentResponse, error) {
	result, err := ps.paymentRepo.CreatePaymentTx(ctx, repository.CreatePaymentTxParams{
		BookingID: req.BookingID,
		UserID:    userID,
		Method:    req.Method,
	})
	if err != nil {
		return nil, err
	}

	// Defensive nil check
	if result == nil || result.Payment == nil {
		return nil, errors.New("payment creation failed")
	}

	// Pass copy vào goroutine — tránh closure capture bug
	seats := result.BookingSeats
	go ps.releaseSeatsAsync(userID, seats)

	return ps.toPaymentResponse(result.Payment), nil
}

// ─────────────────────────────────────────
// releaseSeatsAsync — worker pool release Redis seat hold
// ─────────────────────────────────────────

func (s *paymentService) releaseSeatsAsync(userID string, seats []repository.BookingSeatResult) {
	// Timeout 3 giây — tránh goroutine leak nếu Redis chậm
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const maxWorkers = 4
	workerCount := len(seats)
	if workerCount > maxWorkers {
		workerCount = maxWorkers
	}

	// Channel để distribute jobs
	jobs := make(chan repository.BookingSeatResult, len(seats))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for bs := range jobs {
				if err := cache.ReleaseSeatIfOwner(ctx, bs.TripID, bs.SeatID, userID); err != nil {
					global.Logger.Error("failed to release seat hold",
						zap.String("seat_id", bs.SeatID),
						zap.String("trip_id", bs.TripID),
						zap.String("user_id", userID),
						zap.Error(err),
					)
				}
			}
		}()
	}

	// Send jobs
	for _, bs := range seats {
		jobs <- bs
	}
	close(jobs)

	wg.Wait()
}

// ─────────────────────────────────────────
// Helper
// ─────────────────────────────────────────

func (s *paymentService) toPaymentResponse(payment *repository.PaymentResult) *model.PaymentResponse {
	amount := payment.Amount

	paidAt := ""
	if payment.PaidAt.Valid {
		paidAt = payment.PaidAt.Time.Format(timeFormat.DateTimeFormat)
	}

	return &model.PaymentResponse{
		ID:        payment.ID,
		BookingID: payment.BookingID,
		Amount:    amount,
		Method:    payment.Method,
		Status:    int(payment.Status),
		PaidAt:    paidAt,
	}
}
