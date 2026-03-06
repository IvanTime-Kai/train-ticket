-- +goose Up
-- +goose StatementBegin
CREATE TABLE stations (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    code       VARCHAR(10)  NOT NULL UNIQUE,
    city       VARCHAR(100) NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE trains (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE, -- VD: SE1, SE2, SE3
    total_seats INT          NOT NULL,
    status      TINYINT NOT NULL DEFAULT 1,  -- 1: active | 0: inactive
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE seats (
    id          CHAR(36)      NOT NULL PRIMARY KEY,
    train_id    CHAR(36)      NOT NULL,
    seat_number VARCHAR(10)   NOT NULL,        -- VD: A1, A2, B1
    class       VARCHAR(20)   NOT NULL,        -- economy | business | vip
    price       DECIMAL(12,2) NOT NULL,
    created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_seats_train FOREIGN KEY (train_id) REFERENCES trains(id) ON DELETE CASCADE,
    UNIQUE KEY uq_train_seat (train_id, seat_number) -- 1 tàu không có 2 ghế cùng số
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE routes (
    id                     CHAR(36) NOT NULL PRIMARY KEY,
    origin_station_id      CHAR(36) NOT NULL,      -- ga xuất phát
    destination_station_id CHAR(36) NOT NULL,      -- ga đến
    distance_km            INT,                    -- khoảng cách km
    created_at             DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_routes_origin      FOREIGN KEY (origin_station_id)      REFERENCES stations(id),
    CONSTRAINT fk_routes_destination FOREIGN KEY (destination_station_id) REFERENCES stations(id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE trips (
    id             CHAR(36)   NOT NULL PRIMARY KEY,
    train_id       CHAR(36)   NOT NULL,
    route_id       CHAR(36)   NOT NULL,
    departure_time DATETIME   NOT NULL,            -- giờ khởi hành
    arrival_time   DATETIME   NOT NULL,            -- giờ đến
    status         TINYINT NOT NULL DEFAULT 1,  -- 1: scheduled | 2: departed | 3: arrived | 0: cancelled
    created_at     DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_trips_train FOREIGN KEY (train_id) REFERENCES trains(id),
    CONSTRAINT fk_trips_route FOREIGN KEY (route_id) REFERENCES routes(id),
    INDEX idx_trips_departure (departure_time), -- tìm kiếm theo ngày
    INDEX idx_trips_status (status)             -- filter theo trạng thái
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS trips;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS routes;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS seats;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS trains;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS stations;
-- +goose StatementEnd
