-include .env
export

MYSQL_DSN = $(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DBNAME)

.PHONY: run build start wire migrate-up migrate-down migrate-status migrate-create

wire:
	cd pkg/wire && wire

run:
	go run ./cmd/server/main.go

build:
	go build -o bin/user-service ./cmd/server/main.go

start: build
	./bin/user-service

migrate-up:
	goose -dir db/migrations mysql "$(MYSQL_DSN)" up

migrate-down:
	goose -dir db/migrations mysql "$(MYSQL_DSN)" down

migrate-status:
	goose -dir db/migrations mysql "$(MYSQL_DSN)" status

migrate-create:
	goose -dir db/migrations create $(name) sql