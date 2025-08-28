DB_URL=$(shell grep DB_URL .env | cut -d '=' -f2-)
MIGRATIONS_DIR=./migrations

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) status

test:
	go test ./... -v -cover