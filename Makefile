.PHONY: user-service-run user-service-build user-service-sqlgen train-service-run train-service-build train-service-sqlgen booking-service-run booking-service-build booking-service-sqlgen \
	docker-build-user-service docker-build-train-service docker-build-booking-service docker-up docker-down \
	docker-logs-user-service docker-logs-train-service docker-logs-booking-service

# User service
user-service-run:
	$(MAKE) -C user-service run

user-service-build:
	$(MAKE) -C user-service build

user-service-sqlgen:
	$(MAKE) -C user-service sqlgen

# Train service
train-service-run:
	$(MAKE) -C train-service run

train-service-build:
	$(MAKE) -C train-service build

train-service-sqlgen:
	$(MAKE) -C train-service sqlgen

# Booking service
booking-service-run:
	$(MAKE) -C booking-service run

booking-service-build:
	$(MAKE) -C booking-service build

booking-service-sqlgen:
	$(MAKE) -C booking-service sqlgen


docker-build-user-service:
	@echo "🐳 Building Docker image..."
	@docker compose build user-service

docker-build-train-service:
	@echo "🐳 Building Docker image..."
	@docker compose build train-service

docker-build-booking-service:
	@echo "🐳 Building Docker image..."
	@docker compose build booking-service

docker-up:
	@echo "🐳 Starting all services..."
	@docker compose up -d
docker-down:
	@echo "🐳 Stopping all services..."
	@docker compose down

docker-logs-user-service:
	@echo "📋 Showing logs..."
	@docker compose logs -f user-service

docker-logs-train-service:
	@echo "📋 Showing logs..."
	@docker compose logs -f train-service

docker-logs-booking-service:
	@echo "📋 Showing logs..."
	@docker compose logs -f booking-service