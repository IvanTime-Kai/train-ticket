-include .env
export

MYSQL_DSN = $(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DBNAME)

.PHONY: run build start wire sqlgen migrate-up migrate-down migrate-status migrate-create

sqlgen:
	sqlc generate

wire:
	cd pkg/wire && wire

run:
	go run ./cmd/server/main.go

build:
	go build -o bin/user-service ./cmd/server/main.go

start: build
	./bin/user-service

swag:
	swag init -g cmd/server/main.go

migrate-up:
	goose -dir db/migrations mysql "$(MYSQL_DSN)" up

migrate-down:
	goose -dir db/migrations mysql "$(MYSQL_DSN)" down

migrate-status:
	goose -dir db/migrations mysql "$(MYSQL_DSN)" status

migrate-create:
	goose -dir db/migrations create $(name) sql