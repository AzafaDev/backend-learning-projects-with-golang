include .env
export

.PHONY: run build up down logs migrate-up migrate-down psql clean test

run:
	go run ./main.go

build:
	go build -o bin/server ./...

up:
	docker compose up -d

down:
	docker compose down

up-build:
	docker compose up -d --build

logs:
	docker compose logs -f app

logs-db:
	docker compose logs -f postgres

psql:
	docker compose exec postgres psql -U $(DB_USER) -d $(DB_NAME)

migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

test:
	go test ./... -v

clean:
	docker compose down -v
	rm -rf bin/