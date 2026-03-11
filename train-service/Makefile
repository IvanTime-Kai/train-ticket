-include .env
export

MYSQL_DSN = $(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DBNAME)

.PHONY: run wire swag build start sqlgen migrate-up migrate-down migrate-status migrate-create

sqlgen:
	sqlc generate

run:
	go run ./cmd/server/main.go

wire:
	cd pkg/wire && wire

swag:
	swag init -g cmd/server/main.go --parseDependency --parseInternal --quiet

build:
	go build -o bin/train-service ./cmd/server/main.go

start: build
	./bin/train-service

migrate-up:
	goose -dir db/migrations mysql "$(MYSQL_DSN)" up

migrate-down:
	goose -dir db/migrations mysql "$(MYSQL_DSN)" down

migrate-status:
	goose -dir db/migrations mysql "$(MYSQL_DSN)" status

migrate-create:
	goose -dir db/migrations create $(name) sql