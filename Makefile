# ANSI color codes
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m
COLOR_RED=\033[31m

# Variáveis
MAIN_PATH=cmd/main.go
ENV ?= local
ENV_FILE = .env.$(ENV)
DB_VOLUME_NAME=postgres_data_security

.PHONY: help run install db-start db-stop db-clean db-migrate db-reset build test test-verbose test-coverage test-coverage-html test-domain test-services test-watch dev setup-env


help:
	@echo ""
	@echo "  $(COLOR_YELLOW)Available targets:$(COLOR_RESET)"
	@echo "  $(COLOR_BLUE)Local Development:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)install$(COLOR_RESET)		- Install dependencies"
	@echo "  $(COLOR_GREEN)run$(COLOR_RESET)			- Run development server"
	@echo "  $(COLOR_GREEN)dev$(COLOR_RESET)			- Setup and run complete development environment"
	@echo "  $(COLOR_GREEN)setup-env$(COLOR_RESET)		- Create .env from .env.example"
	@echo ""
	@echo "  $(COLOR_BLUE)Database Management:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)db-start$(COLOR_RESET)		- Start Postgres container"
	@echo "  $(COLOR_GREEN)db-stop$(COLOR_RESET)		- Stop and remove database container"
	@echo "  $(COLOR_GREEN)db-clean$(COLOR_RESET)		- Clean database data (remove container + volume)"
	@echo "  $(COLOR_GREEN)db-migrate$(COLOR_RESET)		- Run database migrations"
	@echo "  $(COLOR_GREEN)db-reset$(COLOR_RESET)		- Reset database (clean + start + migrate)"
	@echo ""
	@echo "  $(COLOR_BLUE)Build & Test:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)build$(COLOR_RESET)		- Build the application"
	@echo "  $(COLOR_GREEN)test$(COLOR_RESET)			- Run all tests"
	@echo "  $(COLOR_GREEN)test-verbose$(COLOR_RESET)		- Run tests with verbose output"
	@echo "  $(COLOR_GREEN)test-coverage$(COLOR_RESET)		- Run tests with coverage report"
	@echo "  $(COLOR_GREEN)test-coverage-html$(COLOR_RESET)	- Generate HTML coverage report"
	@echo "  $(COLOR_GREEN)test-domain$(COLOR_RESET)		- Run domain layer tests only"
	@echo "  $(COLOR_GREEN)test-services$(COLOR_RESET)		- Run services layer tests only"
	@echo "  $(COLOR_GREEN)test-watch$(COLOR_RESET)		- Run tests in watch mode"
	@echo ""

install:
	@echo "$(COLOR_YELLOW)Installing Go dependencies...$(COLOR_RESET)"
	go mod download
	go mod tidy
	@echo "$(COLOR_GREEN)✅ Dependencies installed successfully!$(COLOR_RESET)"

run:
	@echo "$(COLOR_YELLOW)Starting development server with $(ENV) environment...$(COLOR_RESET)"
	@if [ "$(ENV)" = "local" ]; then \
		if [ ! -f .env ]; then \
			if [ -f .env.example ]; then \
				echo "$(COLOR_YELLOW)⚠️  .env not found, creating from .env.example...$(COLOR_RESET)"; \
				cp .env.example .env; \
			else \
				echo "$(COLOR_RED)❌ Neither .env nor .env.example found!$(COLOR_RESET)"; \
				exit 1; \
			fi; \
		fi; \
		echo "$(COLOR_BLUE)Using .env file$(COLOR_RESET)"; \
	else \
		if [ ! -f $(ENV_FILE) ]; then \
			echo "$(COLOR_YELLOW)⚠️  Environment file $(ENV_FILE) not found, creating from .env.example...$(COLOR_RESET)"; \
			cp .env.example $(ENV_FILE); \
		fi; \
		echo "$(COLOR_BLUE)Loading environment from: $(ENV_FILE)$(COLOR_RESET)"; \
		cp $(ENV_FILE) .env; \
	fi
	go run $(MAIN_PATH)

db-start:
	@echo "$(COLOR_YELLOW)Starting Postgres container for $(ENV) environment...$(COLOR_RESET)"
	@if [ "$(ENV)" = "local" ]; then \
		if [ ! -f .env ]; then \
			if [ -f .env.example ]; then \
				echo "$(COLOR_YELLOW)⚠️  .env not found, creating from .env.example...$(COLOR_RESET)"; \
				cp .env.example .env; \
			else \
				echo "$(COLOR_RED)❌ Neither .env nor .env.example found!$(COLOR_RESET)"; \
				exit 1; \
			fi; \
		fi; \
		echo "$(COLOR_BLUE)Using .env file$(COLOR_RESET)"; \
	else \
		if [ ! -f $(ENV_FILE) ]; then \
			echo "$(COLOR_YELLOW)⚠️  Environment file $(ENV_FILE) not found, creating from .env.example...$(COLOR_RESET)"; \
			cp .env.example $(ENV_FILE); \
		fi; \
		echo "$(COLOR_BLUE)Loading environment from: $(ENV_FILE)$(COLOR_RESET)"; \
		cp $(ENV_FILE) .env; \
	fi
	docker compose --env-file .env up postgres -d
	@echo "$(COLOR_GREEN)✅ Database container started!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Database: $$(grep DB_NAME .env | cut -d'=' -f2) on localhost:$$(grep DB_PORT .env | cut -d'=' -f2)$(COLOR_RESET)"

db-stop:
	@echo "$(COLOR_YELLOW)Stopping and removing database container...$(COLOR_RESET)"
	docker compose down postgres
	@echo "$(COLOR_GREEN)✅ Database container removed!$(COLOR_RESET)"

db-clean:
	@echo "$(COLOR_YELLOW)Cleaning database data...$(COLOR_RESET)"
	docker compose down postgres
	docker volume rm $(DB_VOLUME_NAME) 2>/dev/null || true
	@echo "$(COLOR_GREEN)✅ Database data cleaned!$(COLOR_RESET)"

db-migrate:
	@echo "$(COLOR_YELLOW)Running database migrations...$(COLOR_RESET)"
	@if [ ! -f .env ]; then \
		if [ -f .env.example ]; then \
			echo "$(COLOR_YELLOW)⚠️  .env not found, creating from .env.example...$(COLOR_RESET)"; \
			cp .env.example .env; \
		else \
			echo "$(COLOR_RED)❌ Neither .env nor .env.example found!$(COLOR_RESET)"; \
			exit 1; \
		fi; \
	fi
	@echo "$(COLOR_BLUE)Waiting for database to be ready...$(COLOR_RESET)"
	@until docker exec security-postgres pg_isready -U $$(grep DB_USER .env | cut -d'=' -f2) -d $$(grep DB_NAME .env | cut -d'=' -f2) > /dev/null 2>&1; do \
		echo "$(COLOR_YELLOW)⏳ Waiting for database...$(COLOR_RESET)"; \
		sleep 2; \
	done
	@echo "$(COLOR_BLUE)Running SQL migrations...$(COLOR_RESET)"
	docker exec -i security-postgres psql -U $$(grep DB_USER .env | cut -d'=' -f2) -d $$(grep DB_NAME .env | cut -d'=' -f2) < scripts/init.sql
	@echo "$(COLOR_GREEN)✅ Database migrations completed!$(COLOR_RESET)"

db-reset: db-clean db-start db-migrate
	@echo "$(COLOR_GREEN)✅ Database reset completed!$(COLOR_RESET)"

build:
	@echo "$(COLOR_YELLOW)Building application...$(COLOR_RESET)"
	go build -o bin/security-backend $(MAIN_PATH)
	@echo "$(COLOR_GREEN)✅ Build completed! Binary: bin/security-backend$(COLOR_RESET)"

test:
	@echo "$(COLOR_YELLOW)Running tests...$(COLOR_RESET)"
	go test ./...
	@echo "$(COLOR_GREEN)✅ Tests completed!$(COLOR_RESET)"

test-verbose:
	@echo "$(COLOR_YELLOW)Running tests with verbose output...$(COLOR_RESET)"
	go test -v ./...
	@echo "$(COLOR_GREEN)✅ Verbose tests completed!$(COLOR_RESET)"

test-coverage:
	@echo "$(COLOR_YELLOW)Running tests with coverage...$(COLOR_RESET)"
	go test -cover ./...
	@echo "$(COLOR_GREEN)✅ Coverage tests completed!$(COLOR_RESET)"

test-coverage-html:
	@echo "$(COLOR_YELLOW)Generating HTML coverage report...$(COLOR_RESET)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✅ Coverage report generated: coverage.html$(COLOR_RESET)"

test-domain:
	@echo "$(COLOR_YELLOW)Running domain tests...$(COLOR_RESET)"
	go test -v ./internal/core/domain/
	@echo "$(COLOR_GREEN)✅ Domain tests completed!$(COLOR_RESET)"

test-services:
	@echo "$(COLOR_YELLOW)Running services tests...$(COLOR_RESET)"
	go test -v ./internal/core/services/
	@echo "$(COLOR_GREEN)✅ Services tests completed!$(COLOR_RESET)"

test-watch:
	@echo "$(COLOR_YELLOW)Running tests in watch mode...$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Press Ctrl+C to stop$(COLOR_RESET)"
	@while true; do \
		go test ./...; \
		sleep 2; \
	done

dev: install setup-env db-start db-migrate run
	@echo "$(COLOR_GREEN)✅ Development environment ready!$(COLOR_RESET)"

setup-env:
	@echo "$(COLOR_YELLOW)Setting up environment...$(COLOR_RESET)"
	@if [ ! -f .env ]; then \
		if [ -f .env.example ]; then \
			echo "$(COLOR_BLUE)Creating .env from .env.example...$(COLOR_RESET)"; \
			cp .env.example .env; \
			echo "$(COLOR_GREEN)✅ .env file created!$(COLOR_RESET)"; \
		else \
			echo "$(COLOR_RED)❌ .env.example not found!$(COLOR_RESET)"; \
			exit 1; \
		fi; \
	else \
		echo "$(COLOR_BLUE).env file already exists$(COLOR_RESET)"; \
	fi