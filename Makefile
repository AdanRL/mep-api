DEFAULT_GOAL := build

run:
	@echo "🦦🏃‍♂️ Starting Project..."
	@go run ./cmd/api/main.go

build:
	@echo "🛠 Building Project..."
	@go build ./cmd/api/main.go

db-up:
	@echo "🍃🔓 Starting database..."

db-down:
	@echo "🍃🔒 Closing database..."

db-down-rm:
	@echo "🍃🔒 Closing database and removing directory..."

api-dev-up:
	@echo "🐳 Running docker backend container..."
	@docker compose up --build

api-dev-down:
	@echo "🐳 Stopping docker backend container..."
	@docker compose down

api-test-up:
	@echo "🐳 Running docker backend container..."
	@docker compose -f docker-compose-test.yml up --build --exit-code-from api-test

api-test-down:
	@echo "🐳 Stopping docker backend container..."
	@docker compose down

api-pro-up:
	@echo "🐳 Running docker backend container..."
	@docker compose -f docker-compose-prod.yml up -d --build

api-pro-down:
	@echo "🐳 Stopping docker backend container..."
	@docker compose down
