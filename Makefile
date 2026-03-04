.PHONY: user-service-run user-service-build user-service-sqlgen

# User service
user-service-run:
	$(MAKE) -C user-service run

user-service-build:
	$(MAKE) -C user-service build

user-service-sqlgen:
	$(MAKE) -C user-service sqlgen

docker-build:
	@echo "🐳 Building Docker image..."
	@docker compose build user-service

docker-up:
	@echo "🐳 Starting all services..."
	@docker compose up -d

docker-down:
	@echo "🐳 Stopping all services..."
	@docker compose down

docker-logs:
	@echo "📋 Showing logs..."
	@docker compose logs -f user-service