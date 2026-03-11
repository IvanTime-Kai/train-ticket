-- +goose Up
-- +goose StatementBegin
CREATE TABLE bookings (
    id          CHAR(36)      NOT NULL PRIMARY KEY,
    user_id     CHAR(36)      NOT NULL,                  -- từ user-service
    trip_id     CHAR(36)      NOT NULL,                  -- từ train-service
    total_price DECIMAL(12,2) NOT NULL,
    status      TINYINT       NOT NULL DEFAULT 1,        -- 1: pending | 2: confirmed | 3: expired | 0: cancelled
    expires_at  DATETIME      NOT NULL,                  -- thời gian hết hạn thanh toán (15 phút)
    created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_bookings_user   (user_id),
    INDEX idx_bookings_trip   (trip_id),
    INDEX idx_bookings_status (status)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE booking_seats (
    id          CHAR(36)      NOT NULL PRIMARY KEY,
    booking_id  CHAR(36)      NOT NULL,
    seat_id     CHAR(36)      NOT NULL,                  -- từ train-service
    trip_id     CHAR(36)      NOT NULL,
    seat_number VARCHAR(10)   NOT NULL,
    class       VARCHAR(20)   NOT NULL,
    price       DECIMAL(12,2) NOT NULL,
    created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_booking_seats_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    UNIQUE KEY uq_trip_seat (trip_id, seat_id)           -- 1 ghế chỉ được book 1 lần trong 1 chuyến
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE payments (
    id         CHAR(36)      NOT NULL PRIMARY KEY,
    booking_id CHAR(36)      NOT NULL,
    amount     DECIMAL(12,2) NOT NULL,
    method     VARCHAR(20)   NOT NULL DEFAULT 'mock',   -- mock | vnpay | momo
    status     TINYINT       NOT NULL DEFAULT 1,        -- 1: pending | 2: success | 0: failed
    paid_at    DATETIME      NULL,
    created_at DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_payments_booking FOREIGN KEY (booking_id) REFERENCES bookings(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS booking_seats;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS bookings;
-- +goose StatementEnd