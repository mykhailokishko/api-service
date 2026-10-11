.PHONY: all build run test bench lint cover docker-build docker-up docker-down clean

APP_NAME = card-validator
BIN_DIR = bin

all: lint test build

build:
	@echo "==> Building binary..."
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) ./cmd/api

run:
	@echo "==> Running service..."
	go run ./cmd/api

test:
	@echo "==> Running unit and integration tests with race detector..."
	go test -v -race ./...

cover:
	@echo "==> Running tests with coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

bench:
	@echo "==> Running performance benchmarks..."
	go test -v -bench=. -benchmem ./pkg/luhn/...

lint:
	@echo "==> Running linter..."
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not installed, running go vet" && go vet ./...)
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run ./... || true

docker-build:
	@echo "==> Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-up:
	@echo "==> Starting containers via docker compose..."
	docker compose up -d --build

docker-down:
	@echo "==> Stopping containers..."
	docker compose down

clean:
	@echo "==> Cleaning up artifacts..."
	rm -rf $(BIN_DIR) coverage.out
