.PHONY: build test lint vet run migrate docker-build docker-up docker-down clean

APP_NAME := currencyflow
BUILD_DIR := ./bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server/

test:
	go test -v -race -coverprofile=coverage.out ./...

test-cover: test
	go tool cover -html=coverage.out -o coverage.html

lint:
	test -z "$$(gofmt -l .)"
	go vet ./...

run:
	go run ./cmd/server/

migrate:
	go run ./cmd/server/ --migrate

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

.DEFAULT_GOAL := build
