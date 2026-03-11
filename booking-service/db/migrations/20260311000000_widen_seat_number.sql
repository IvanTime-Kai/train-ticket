-- +goose Up
-- +goose StatementBegin
ALTER TABLE booking_seats
MODIFY COLUMN seat_number VARCHAR(36) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE booking_seats
MODIFY COLUMN seat_number VARCHAR(10) NOT NULL;
-- +goose StatementEnd
