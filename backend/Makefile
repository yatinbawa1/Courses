# -include .env

# Variables
DB_USER=postgres
DB_PASSWORD=COURSES
DB_HOST=localhost
DB_PORT=5432
DB_NAME=stage

# Construct the migration database connection URL
DATABASE_URL="postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
MIGRATIONS_PATH=internal/database/migrations

.PHONY: run migrate-up migrate-down migrate-force delete-all-users

#Run App
run:
	go run cmd/main.go

# Run all up migrations
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) up

# Rollback the last migration
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) down 1

# Force clear a dirty migration state (usage: make migrate-force v=1)
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) force $(v)

delete-all-users:
	docker exec -i courses-psql psql $(DATABASE_URL) < internal/database/migrations/000001_delete_all_user.sql

export