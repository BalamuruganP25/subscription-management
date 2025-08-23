# 🛠️ Makefile for subscription-management Project

SERVICE_NAME = subscription-management

export GO111MODULE = on

.PHONY: run build stop dep test

## 🔄 Run the full Docker stack
run: build
	@echo "🚀 Starting $(SERVICE_NAME)..."
	@docker compose up

## 🏗️ Build Docker containers
build:
	@echo "🔧 Building Docker images..."
	@docker compose build

## run unit tests
test:
	@echo "🧪 Running unit tests..."
	go test ./... 

## 🛑 Stop and remove containers
stop:
	@echo "🧹 Stopping services..."
	@docker compose down

## 📦 Install and tidy Go dependencies
dep:
	@echo "📦 Tidying and vendoring Go dependencies..."
	@go mod tidy
	@go mod vendor