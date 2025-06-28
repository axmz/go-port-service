APP_NAME := go-port-service
PKG := ./...
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)
GO := go

.PHONY: dc all build run clean test check lint fmt vet tidy deps cover

all: build

build:
    $(GO) build -o $(BIN) .

build-race: ## Build the app binary with race detector for CI
	$(GO) build -race -o $(BIN) .

run: build
    ./$(BIN)

dc:
	docker-compose up  --remove-orphans --build

clean:
    rm -rf $(BIN_DIR) $(APP_NAME)

# TESTING
test:
    $(GO) test -race -v $(PKG) # -count 1000 -failfast

cover:
    $(GO) test -coverprofile=coverage.out $(PKG)
    $(GO) tool cover -func=coverage.out

cover-html:
	$(GO) tool cover -html=coverage.out

# CODE QUALITY
check: 
	fmt vet lint test

lint:
    golangci-lint run # go tool golanci-lint run

fmt:
    $(GO) fmt $(PKG)

vet:
    $(GO) vet $(PKG)

# MOD
tidy:
    $(GO) mod tidy

deps:
    $(GO) mod download