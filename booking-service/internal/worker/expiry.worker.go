package worker

import (
	"context"
	"time"

	"github.com/leminhthai/train-ticket/booking-service/internal/repository"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type ExpiryWorker struct {
	bookingRepo repository.BookingRepository
	logger      *zap.Logger
	cron        *cron.Cron
}

func NewExpiryWorker(bookingRepo repository.BookingRepository, logger *zap.Logger) *ExpiryWorker {
	return &ExpiryWorker{
		bookingRepo: bookingRepo,
		logger:      logger,
		cron: cron.New(
			cron.WithLogger(cron.VerbosePrintfLogger(nil)),
			cron.WithChain(
				cron.Recover(cron.DefaultLogger),
			),
		),
	}
}

func (w *ExpiryWorker) Start() {
	w.cron.AddFunc("* * * * *", w.expireBookings)

	w.cron.Start()

	w.logger.Info("ExpiryWorker started")

	// Chạy ngay lần đầu khi khởi động
	go w.expireBookings()
}

func (w *ExpiryWorker) Stop() {
	ctx := w.cron.Stop()

	select {
	case <-ctx.Done():
		w.logger.Info("ExpiryWorker stopped gracefully")
	case <-time.After(10 * time.Second):
		w.logger.Warn("ExpiryWorker stop timeout")
	}
}

func (w *ExpiryWorker) expireBookings() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	expired, err := w.bookingRepo.GetExpiredBookings(ctx)

	if err != nil {
		w.logger.Error("ExpiryWorker: GetExpiredBookings failed", zap.Error(err))
		return
	}

	if len(expired) == 0 {
		return
	}

	if err := w.bookingRepo.BulkUpdateExpiredBookings(ctx); err != nil {
		w.logger.Error("ExpiryWorker: BulkUpdateExpiredBookings failed", zap.Error(err))
		return
	}

	w.logger.Info("ExpiryWorker: bookings expired",
		zap.Int("count", len(expired)),
	)
}
