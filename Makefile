.PHONY: up down build-up lint-plan-test lint readme-lint validate-sam test build-lambdas

# 🚀 Project-specific settings
APP_NAME := equilibria
DOCKER_COMPOSE_FILE := ./docker-compose.yml
MIGRATIONS_DIR := ./migrations

run:
	@echo "🚀 Starting Local API..."
	cd lambdas && sam build && sam local start-api

# Docker Compose Commands
docker-up:
	@echo "🚀 Starting Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "🛑 Stopping Docker Compose..."
	docker-compose down

docker-build:
	@echo "🔨 Building and Starting Docker Compose..."
	docker-compose up --build -d

docker: generate-sql-init docker-build

# Terraform Commands
terraform-plan:
	@echo "📝 Running Terraform Plan..."
	cd terraform && terraform plan

# Linting Commands
lint: go-lint readme-lint

go-lint:
	@echo "🔍 Linting Go files..."
	cd lambdas && golangci-lint run ./...

readme-lint:
	@echo "📖 Linting README files..."
	find . -name '*.md' -exec markdownlint {} +

# SAM Template Validation
validate-sam:
	@echo "🔎 Validating SAM Template..."
	sam validate --template-file templates/template.yaml

# Combined Lint, Validate, and Plan Command
lint-plan-test: lint readme-lint validate-sam terraform-plan

# Testing Commands
test:
	@echo "🧪 Running all tests..."
	source .env && go test -json -v ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

convey:
	@echo "🧪 Conveying tests in browser..."
	source .env && goconvey -excludedDirs=vendor

# Build all sms Lambda Functions
build: go-lint build-authorizer build-login build-receive-sms build-send-sms build-status-sms build-manage-user build-signup-otp build-nudger-sms build-factfinder

# Build authorizer lambda function
build-authorizer:
	@echo "🛠 Building Authorizer lambda..."
	cd lambdas/authorizer && GOOS=linux GOARCH=amd64 go build -o main && \
	cp ../../build/bootstrap . && \
	zip authorizer.zip main bootstrap && \
	rm main bootstrap && mv authorizer.zip ../../build

# Build login lambda function
build-login:
	@echo "🛠 Building Login lambda..."
	cd lambdas/login && GOOS=linux GOARCH=amd64 go build -o main && \
	cp ../../build/bootstrap . && \
	zip login.zip main bootstrap && \
	rm main bootstrap && mv login.zip ../../build

# Build FactFinder lambda Functions
build-factfinder:
	@echo "🛠 Building FactFinder lambda..."
	cd lambdas/factfinder && GOOS=linux GOARCH=amd64 go build -o main && \
	cp ../../build/bootstrap . && \
	zip factfinder.zip main bootstrap && \
	rm main bootstrap && mv factfinder.zip ../../build

# Build status lambda Functions
build-status-sms:
	@echo "🛠 Building SMS Status lambda..."
	cd lambdas/status_sms && GOOS=linux GOARCH=amd64 go build -o main && \
	cp ../../build/bootstrap . && \
	zip status_sms.zip main bootstrap && \
	rm main bootstrap && mv status_sms.zip ../../build

# Build all sms Lambda Functions
build-receive-sms:
	@echo "🛠 Building SMS Receiver lambda..."
	cd lambdas/receive_sms && GOOS=linux GOARCH=amd64 go build -o main && \
	cp ../../build/bootstrap . && \
	zip receive_sms.zip main bootstrap && \
	rm main bootstrap && mv receive_sms.zip ../../build

build-send-sms:
	@echo "🛠 Building SMS Sender lambda..."
	cd lambdas/send_sms && GOOS=linux GOARCH=amd64 go build -o main && \
	cp ../../build/bootstrap . && \
	zip send_sms.zip main bootstrap && \
	rm main bootstrap && mv send_sms.zip ../../build

build-signup-otp:
	@echo "🛠 Building SMS OTP lambda..."
	cd lambdas/signup_otp && GOOS=linux GOARCH=amd64 go build -o main && \
	cp ../../build/bootstrap . && \
	zip signup_otp.zip main bootstrap && \
	rm main bootstrap && mv signup_otp.zip ../../build

build-manage-user:
	@echo "🛠 Building User Mgmt lambda..."
	cd lambdas/manage_user && GOOS=linux GOARCH=amd64 go build -o main && \
	cp ../../build/bootstrap . && \
	zip manage_user.zip main bootstrap && \
	rm main bootstrap && mv manage_user.zip ../../build

build-nudger-sms:
	@echo "🛠 Building Nudger lambda..."
	cd lambdas/nudge_sms && GOOS=linux GOARCH=amd64 go build -o main && \
	cp ../../build/bootstrap . && \
	zip nudge_sms.zip main bootstrap && \
	rm main bootstrap && mv nudge_sms.zip ../../build

# 🗃️ Perform database migrations
migrate:
	@echo "🗃️ Performing database migrations..."
	goose -dir ${MIGRATIONS_DIR} mysql "${DATABASE_USER}:${MYSQL_ROOT_PASSWORD}@tcp(${DATABASE_HOST})/${DATABASE_NAME}" up

# 🗃️ Check status of database migrations
db-status:
	@echo "🗃️ Checking status of database migrations..."
	goose -dir ${MIGRATIONS_DIR} mysql "${DATABASE_USER}:${MYSQL_ROOT_PASSWORD}@tcp(${DATABASE_HOST})/${DATABASE_NAME}" status

rollback:
	@echo "🗃️ Performing database migrations..."
	# Add your database migration command here
	# Example: goose -dir $(MIGRATIONS_DIR) mysql "user:password@/dbname" up
	goose -dir ${MIGRATIONS_DIR} mysql "${DATABASE_USER}:${MYSQL_ROOT_PASSWORD}@tcp(${DATABASE_HOST})/${DATABASE_NAME}" down

# 🗑️ Clear the database by rolling back all migrations
clear-database:
	@echo "🗑️ Clearing the entire database..."
	@while goose -dir ${MIGRATIONS_DIR} mysql "${DATABASE_USER}:${MYSQL_ROOT_PASSWORD}@tcp(${DATABASE_HOST})/${DATABASE_NAME}" down && [ $$? -eq 0 ]; do \
		echo "Rolling back migration..."; \
	done

generate-sql-init:
	@echo "📝 Generating SQL init file..."
	@envsubst < ./init.sql > ./build/init.sql
