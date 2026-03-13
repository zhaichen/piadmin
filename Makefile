APP_NAME := piadmin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

.PHONY: all frontend build build-arm7 build-arm64 build-all dev-backend dev-frontend clean test

all: frontend build

frontend:
	cd web/frontend && npm ci && npm run build

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) ./cmd/piadmin

build-arm7: frontend
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-linux-arm7 ./cmd/piadmin

build-arm64: frontend
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-linux-arm64 ./cmd/piadmin

build-all: frontend build-arm7 build-arm64

dev-backend:
	go run ./cmd/piadmin -config configs/piadmin.example.yaml

dev-frontend:
	cd web/frontend && npm run dev

clean:
	rm -rf bin/ web/frontend/dist web/frontend/node_modules

test:
	go test ./... -v -race
