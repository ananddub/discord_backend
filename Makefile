# Variables
DB_URL := postgresql://postgres:postgres@localhost:5432/discord?sslmode=disable
GOOSE_DRIVER := postgres
MIGRATIONS_DIR := ./sql/migrations

# Development Commands
.PHONY: dev watch build run

dev: build
	./bin/discord

watch:
	air

# Build Commands
build:
	go build -o bin/discord ./cmd/grpc
	@echo "✅ Build completed"

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/discord-linux ./cmd/grpc
	@echo "✅ Linux build completed"

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/discord.exe ./cmd/grpc
	@echo "✅ Windows build completed"

build-all: build build-linux build-windows
	@echo "✅ All builds completed"

run: build
	./bin/discord

run-dev:
	air

# Proto Generation
proto-gen:
	buf generate
	@echo "✅ Protobuf files generated"

proto-lint:
	buf lint
	@echo "✅ Proto files linted"

proto-format:
	buf format -w
	@echo "✅ Proto files formatted"

proto-breaking:
	buf breaking --against '.git#branch=main'

# SQLC Commands
sqlc-gen:
	sqlc generate
	@echo "✅ SQLC files generated"

sqlc-compile:
	sqlc compile
	@echo "✅ SQL queries compiled"

# Database Commands
db-create:
	psql -U postgres -c "CREATE DATABASE discord;"
	@echo "✅ Database created"

db-drop:
	psql -U postgres -c "DROP DATABASE IF EXISTS discord;"
	@echo "✅ Database dropped"

db-reset: db-drop db-create db-migrate
	@echo "✅ Database reset completed"

db-migrate:
	goose -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DB_URL)" up
	@echo "✅ Migrations completed"

db-migrate-down:
	goose -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DB_URL)" down
	@echo "✅ Migration rolled back"

db-migrate-status:
	goose -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) "$(DB_URL)" status

db-migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql
	@echo "✅ Migration $(NAME) created"

db-seed:
	go run cmd/seed/main.go
	@echo "✅ Database seeded"

db-psql:
	psql "$(DB_URL)"

# Code Generation
gen-all: proto-gen sqlc-gen
	@echo "✅ All code generated"

# Testing
test:
	go test ./... -v

test-unit:
	go test ./... -v -short

test-integration:
	go test ./... -v -run Integration

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Code Quality
lint:
	golangci-lint run
	buf lint

fmt:
	gofmt -s -w .
	goimports -w .
	@echo "✅ Code formatted"

vet:
	go vet ./...

# Dependencies
deps-install:
	go mod download
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ All dependencies installed"

deps-tidy:
	go mod tidy
	@echo "✅ Dependencies tidied"

deps-update:
	go get -u ./...
	go mod tidy
	@echo "✅ Dependencies updated"

# Docker Commands
docker-build:
	docker build -t discord-backend:latest .
	@echo "✅ Docker image built"

docker-run:
	docker run -p 50051:50051 discord-backend:latest

docker-compose-up:
	docker-compose up -d
	@echo "✅ Services started"

docker-compose-down:
	docker-compose down
	@echo "✅ Services stopped"

# Cleanup
clean:
	rm -rf bin/
	rm -rf gen/proto/
	find . -name "*.pb.go" -delete
	rm -f coverage.out coverage.html
	@echo "✅ Cleaned"

clean-all: clean
	go clean -cache -testcache -modcache
	@echo "✅ Deep clean completed"

# Project Setup
setup: deps-install db-create db-migrate gen-all
	@echo "✅ Project setup completed"

# Production
prod-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o bin/discord ./cmd/grpc
	@echo "✅ Production build completed"

prod-deploy:
	@echo "🚀 Deploying to production..."

rs-client:
	protoc -I /mnt/d/Devloper/test/test_server/proto \
		--go_out=gen --go_opt=paths=source_relative \
		--go-grpc_out=gen --go-grpc_opt=paths=source_relative \
		$$(find /mnt/d/Devloper/test/test_server/proto -name "*.proto")

# Help
help:
	@echo "Available targets:"
	@echo "  dev              - Run in development mode"
	@echo "  build            - Build the application"
	@echo "  test             - Run tests"
	@echo "  proto-gen        - Generate protobuf files"
	@echo "  sqlc-gen         - Generate SQLC files"
	@echo "  db-migrate       - Run database migrations"
	@echo "  clean            - Clean generated files"
	@echo "  setup            - Initial project setup"
