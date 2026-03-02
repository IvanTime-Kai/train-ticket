.PHONY: run sqlgen build start

run:
	go run ./cmd/server/main.go

sqlgen:
	sqlc generate

build:
	go build -o bin/user-service ./cmd/server/main.go

start: build
	./bin/user-service