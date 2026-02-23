.PHONY: generate up down run

# Generates Go code from SQL queries using sqlc via Docker
generate:
	docker run --rm -v "$$(pwd):/src" -w /src sqlc/sqlc generate

# Starts the local infrastructure (Postgres, Redis)
up:
	docker compose up -d

# Stops the local infrastructure
down:
	docker compose down

# Runs the Go application locally
run:
	go run cmd/api/main.go