.PHONY: help build up down restart logs clean test docker-build

help: ## Show help
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compile local binary
	go build -o bin/api ./cmd/product-comparison-api

run: ## Run local
	go run ./cmd/product-comparison-api

test: ## Run tests
	go test -v ./...

docker-build: ## Build Docker image
	docker compose build

up: ## Start all services
	docker compose up -d

up-build: ## Build and start all services
	docker compose up -d --build

down: ## Stop all services
	docker compose down

down-volumes: ## Stop and remove volumes
	docker compose down -v

restart: ## Restart services
	docker compose restart

restart-api: ## Restart only the API
	docker compose restart product-comparison-api

logs: ## View logs of all services
	docker compose logs -f

logs-api: ## View logs only of the API
	docker compose logs -f product-comparison-api

logs-grafana: ## View logs of Grafana
	docker compose logs -f grafana

health: ## Verify health check
	@curl -s http://localhost:8080/api/health-check

test-compare: ## Test comparison endpoint
	@curl -s -X POST http://localhost:8080/api/v1/items/compare \
		-H "Content-Type: application/json" \
		-d '{"ids":["4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d","30106bcd-f425-4dfb-8ef6-055ab4744f6c"]}'

stats: ## View container statistics
	docker stats --no-stream