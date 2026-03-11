package model

// ─────────────────────────────────────────
// Booking Status Constants
// ─────────────────────────────────────────
const (
	BookingStatusPending   = 1 // vừa tạo, chờ thanh toán
	BookingStatusConfirmed = 2 // đã thanh toán
	BookingStatusCancelled = 0 // đã huỷ
	BookingStatusExpired   = 3 // quá 15 phút chưa thanh toán
)

// Payment Status Constants
const (
	PaymentStatusPending = 1 // chờ thanh toán
	PaymentStatusSuccess = 2 // thanh toán thành công
	PaymentStatusFailed  = 0 // thanh toán thất bại
)

// Payment Method Constants
const (
	PaymentMethodMock  = "mock"
	PaymentMethodVNPay = "vnpay"
	PaymentMethodMomo  = "momo"
)

// ─────────────────────────────────────────
// Request DTOs
// ─────────────────────────────────────────

type HoldSeatRequest struct {
	TripID  string   `json:"trip_id"  binding:"required"`
	SeatIDs []string `json:"seat_ids" binding:"required,min=1"`
}

type CreateBookingRequest struct {
	TripID    string   `json:"trip_id"    binding:"required"`
	SeatIDs   []string `json:"seat_ids"   binding:"required,min=1"`
	HoldToken string   `json:"hold_token" binding:"required"`
}

type CancelBookingRequest struct {
	BookingID string `json:"booking_id" binding:"required"`
}

type CreatePaymentRequest struct {
	BookingID string `json:"booking_id" binding:"required"`
	Method    string `json:"method"     binding:"required,oneof=mock vnpay momo"`
}

// ─────────────────────────────────────────
// Response DTOs
// ─────────────────────────────────────────

type SeatInfo struct {
	SeatID     string  `json:"seat_id"`
	SeatNumber string  `json:"seat_number"`
	Class      string  `json:"class"`
	Price      float64 `json:"price"`
}

type HoldSeatResponse struct {
	HoldToken string     `json:"hold_token"`
	TripID    string     `json:"trip_id"`
	Seats     []SeatInfo `json:"seats"`
	ExpiresIn int        `json:"expires_in"` // giây
}

type BookingResponse struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	TripID     string     `json:"trip_id"`
	TotalPrice float64    `json:"total_price"`
	Status     int        `json:"status"`
	ExpiresAt  string     `json:"expires_at"`
	Seats      []SeatInfo `json:"seats"`
	CreatedAt  string     `json:"created_at"`
}

type PaymentResponse struct {
	ID        string  `json:"id"`
	BookingID string  `json:"booking_id"`
	Amount    float64 `json:"amount"`
	Method    string  `json:"method"`
	Status    int     `json:"status"`
	PaidAt    string  `json:"paid_at"`
}