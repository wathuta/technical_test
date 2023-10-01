include .env.dev

DATABASE_URL = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_FOLDER = $(PWD)/internal/platform/migrations


unit-test:
	go test -v -timeout 30s -coverprofile=cover.out -cover ./internal/...
	go tool cover -func=cover.out
	rm cover.out

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

flush-db:
	psql $(DATABASE_URL) -c 'DROP SCHEMA public CASCADE'
	psql $(DATABASE_URL) -c 'CREATE SCHEMA public'

