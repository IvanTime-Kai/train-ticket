package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	db "github.com/leminhthai/train-ticket/booking-service/db/generated"
	"github.com/leminhthai/train-ticket/booking-service/internal/model"
	"github.com/leminhthai/train-ticket/booking-service/internal/repository"
	"github.com/leminhthai/train-ticket/booking-service/internal/utils/cache"
	timeFormat "github.com/leminhthai/train-ticket/booking-service/internal/utils/time"

	grpcClient "github.com/leminhthai/train-ticket/booking-service/internal/grpc"
	proto "github.com/IvanTime-Kai/train-ticket-proto/gen/train"
)

type BookingService interface {
	HoldSeat(ctx context.Context, userID string, req *model.HoldSeatRequest) (*model.HoldSeatResponse, error)
	CreateBooking(ctx context.Context, userID string, req *model.CreateBookingRequest) (*model.BookingResponse, error)
	GetBookingByID(ctx context.Context, userID, bookingID string) (*model.BookingResponse, error)
	GetMyBookings(ctx context.Context, userID string) ([]model.BookingResponse, error)
	CancelBooking(ctx context.Context, userID, bookingID string) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	trainClient *grpcClient.TrainClient
}

func NewBookingService(bookingRepo repository.BookingRepository, trainClient *grpcClient.TrainClient) BookingService {
	return &bookingService{bookingRepo: bookingRepo, trainClient: trainClient}
}

// ─────────────────────────────────────────
// Hold Seat
// ─────────────────────────────────────────

func (s *bookingService) HoldSeat(ctx context.Context, userID string, req *model.HoldSeatRequest) (*model.HoldSeatResponse, error) {
	// 1. Validate seats qua gRPC TRƯỚC — không tốn Redis nếu seat invalid
	protoSeats, err := s.trainClient.ValidateSeats(ctx, req.TripID, req.SeatIDs)
	if err != nil {
		return nil, err
	}
	seats := protoSeatsToModel(protoSeats)

	// 2. Check DB
	bookedSeatIDs, err := s.bookingRepo.AreSeatsBooked(ctx, req.TripID, req.SeatIDs)
	if err != nil {
		return nil, err
	}
	if len(bookedSeatIDs) > 0 {
		return nil, fmt.Errorf("seats already booked: %v", bookedSeatIDs)
	}

	// 3. Hold Redis atomic
	if err := cache.HoldMultipleSeatsAtomic(ctx, req.TripID, req.SeatIDs, userID); err != nil {
		return nil, err
	}

	// 4. Lưu token
	token, err := cache.SaveHoldToken(ctx, userID, req.TripID, req.SeatIDs)
	if err != nil {
		for _, seatID := range req.SeatIDs {
			_ = cache.ReleaseSeatIfOwner(ctx, req.TripID, seatID, userID)
		}
		return nil, err
	}

	return &model.HoldSeatResponse{
		HoldToken: token,
		TripID:    req.TripID,
		Seats:     seats,
		ExpiresIn: 300,
	}, nil
}

// ─────────────────────────────────────────
// Create Booking
// ─────────────────────────────────────────

func (s *bookingService) CreateBooking(ctx context.Context, userID string, req *model.CreateBookingRequest) (*model.BookingResponse, error) {
	// Verify hold token
	holdData, err := cache.GetHoldToken(ctx, req.HoldToken)
	if err != nil {
		return nil, err
	}

	if holdData.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	if holdData.TripID != req.TripID {
		return nil, fmt.Errorf("invalid hold token")
	}

	if len(holdData.SeatIDs) != len(req.SeatIDs) {
		return nil, fmt.Errorf("invalid hold token")
	}
	seatMap := make(map[string]bool)
	for _, id := range holdData.SeatIDs {
		seatMap[id] = true
	}
	for _, id := range req.SeatIDs {
		if !seatMap[id] {
			return nil, fmt.Errorf("invalid hold token")
		}
	}

	// Verify ownership + extend TTL — 1 Lua script cho tất cả ghế
	// Atomic: check owner all + extend all
	if err := cache.ExtendMultipleSeatsIfOwner(ctx, req.TripID, req.SeatIDs, userID); err != nil {
		return nil, err
	}

	// Validate seats + lấy thông tin thật từ train-service qua gRPC
	protoSeats, err := s.trainClient.ValidateSeats(ctx, req.TripID, req.SeatIDs)
	if err != nil {
		return nil, err
	}
	seats := protoSeatsToModel(protoSeats)
	totalPrice := 0.0
	for i := range seats {
		totalPrice += seats[i].Price
	}

	bookingID := uuid.New().String()
	expiresAt := time.Now().Add(15 * time.Minute)

	// DB Transaction — INSERT booking + booking_seats
	if err := s.bookingRepo.CreateBookingTx(ctx, repository.CreateBookingTxParams{
		BookingID:  bookingID,
		UserID:     userID,
		TripID:     req.TripID,
		TotalPrice: totalPrice,
		ExpiresAt:  expiresAt,
		Seats:      seats,
	}); err != nil {
		if isDuplicateError(err) {
			return nil, fmt.Errorf("seat already booked")
		}
		return nil, err
	}

	// Xoá hold token SAU KHI commit thành công
	_ = cache.DeleteHoldToken(ctx, req.HoldToken)

	// Release seat hold — chỉ release khi đúng owner
	for _, seatID := range req.SeatIDs {
		_ = cache.ReleaseSeatIfOwner(ctx, req.TripID, seatID, userID)
	}

	// Query lại từ DB → lấy created_at chính xác
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	return s.toBookingResponse(booking, seats), nil
}

// ─────────────────────────────────────────
// Get Booking By ID
// ─────────────────────────────────────────

func (s *bookingService) GetBookingByID(ctx context.Context, userID, bookingID string) (*model.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.UserID != userID {
		return nil, fmt.Errorf("forbidden")
	}

	seats, err := s.getBookingSeats(ctx, booking)
	if err != nil {
		return nil, err
	}

	return s.toBookingResponse(booking, seats), nil
}

// ─────────────────────────────────────────
// Get My Bookings
// ─────────────────────────────────────────

func (s *bookingService) GetMyBookings(ctx context.Context, userID string) ([]model.BookingResponse, error) {
	bookings, err := s.bookingRepo.GetBookingsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]model.BookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		b := booking
		seats, err := s.getBookingSeats(ctx, &b)
		if err != nil {
			return nil, err
		}
		result = append(result, *s.toBookingResponse(&b, seats))
	}

	return result, nil
}

// ─────────────────────────────────────────
// Cancel Booking
// ─────────────────────────────────────────

func (s *bookingService) CancelBooking(ctx context.Context, userID, bookingID string) error {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking.UserID != userID {
		return fmt.Errorf("forbidden")
	}

	switch int(booking.Status) {
	case model.BookingStatusConfirmed:
		return fmt.Errorf("cannot cancel confirmed booking")
	case model.BookingStatusCancelled:
		return fmt.Errorf("booking already cancelled")
	case model.BookingStatusExpired:
		return fmt.Errorf("booking already expired")
	}

	// Chỉ update DB — không cần release Redis
	// Seat availability dựa vào booking_seats table
	// Redis hold đã bị xoá sau CreateBooking
	return s.bookingRepo.UpdateBookingStatus(ctx, bookingID, int8(model.BookingStatusCancelled))
}

// ─────────────────────────────────────────
// Helper
// ─────────────────────────────────────────

func (s *bookingService) getBookingSeats(ctx context.Context, booking *db.Booking) ([]model.SeatInfo, error) {
	bookingSeats, err := s.bookingRepo.GetBookingSeatsByBookingID(ctx, booking.ID)
	if err != nil {
		return nil, err
	}

	seats := make([]model.SeatInfo, 0, len(bookingSeats))
	for _, bs := range bookingSeats {
		price, _ := strconv.ParseFloat(bs.Price, 64)
		seats = append(seats, model.SeatInfo{
			SeatID:     bs.SeatID,
			SeatNumber: bs.SeatNumber,
			Class:      bs.Class,
			Price:      price,
		})
	}

	return seats, nil
}

func (s *bookingService) toBookingResponse(booking *db.Booking, seats []model.SeatInfo) *model.BookingResponse {
	totalPrice, _ := strconv.ParseFloat(booking.TotalPrice, 64)
	return &model.BookingResponse{
		ID:         booking.ID,
		UserID:     booking.UserID,
		TripID:     booking.TripID,
		TotalPrice: totalPrice,
		Status:     int(booking.Status),
		ExpiresAt:  booking.ExpiresAt.Format(timeFormat.DateTimeFormat),
		Seats:      seats,
		CreatedAt:  booking.CreatedAt.Format(timeFormat.DateTimeFormat),
	}
}

func isDuplicateError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062 // MySQL duplicate entry
	}
	return false
}

func protoSeatsToModel(seats []*proto.SeatInfo) []model.SeatInfo {
	out := make([]model.SeatInfo, 0, len(seats))
	for _, s := range seats {
		if s == nil {
			continue
		}
		out = append(out, model.SeatInfo{
			SeatID:     s.SeatId,
			SeatNumber: s.SeatNumber,
			Class:      s.Class,
			Price:      s.Price,
		})
	}
	return out
}
