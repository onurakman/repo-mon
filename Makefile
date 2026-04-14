GO_BIN := $(HOME)/go/go1.26.2/bin
export PATH := $(GO_BIN):$(HOME)/go/bin:$(PATH)

WAILS_TAGS := webkit2_41
APP_NAME := repo-mon

.PHONY: dev build test clean install-deps doctor

## Development

dev: ## Run in dev mode with hot reload
	wails dev -tags $(WAILS_TAGS)

## Build

build: ## Build production binary
	wails build -tags $(WAILS_TAGS)

## Testing

test: ## Run Go tests
	go test ./internal/... -v

test-git: ## Run git wrapper tests only
	go test ./internal/git/ -v

## Dependencies

install-deps: ## Install all dependencies (Go + frontend)
	go mod tidy
	cd frontend && npm install --legacy-peer-deps

## Maintenance

clean: ## Remove build artifacts
	rm -rf build/bin/*
	rm -rf frontend/dist
	rm -rf frontend/node_modules

doctor: ## Run Wails diagnostics
	wails doctor

generate: ## Regenerate Wails bindings
	wails generate module

## Help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
