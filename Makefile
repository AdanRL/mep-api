DEFAULT_GOAL := build
include .env.dev

COMPOSE_DEV = docker compose -f docker-compose.yml --env-file .env.dev
COMPOSE_PROD = docker compose -f docker-compose-prod.yml --env-file .env.prod

dev-up:
	@$(COMPOSE_DEV) up -d db api

dev-down:
	@$(COMPOSE_DEV) down

migrate-up:
	@$(COMPOSE_DEV) run --rm migrator -path=/migrations -database="$(DATABASE_URL)" up

migrate-down:
	@$(COMPOSE_DEV) run --rm migrator -path=/migrations -database="$(DATABASE_URL)" down 1
