.PHONY: up down logs tidy build test db-reset db-shell

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

db-reset:
	docker compose down -v
	docker compose up -d postgres

db-shell:
	docker compose exec postgres sh -c 'psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"'

tidy:
	cd services/auth && go mod tidy
	cd services/user && go mod tidy
	cd services/subscription && go mod tidy
	cd services/notification && go mod tidy

build:
	cd services/auth && go build ./...
	cd services/user && go build ./...
	cd services/subscription && go build ./...
	cd services/notification && go build ./...

test:
	cd services/auth && go test ./...
	cd services/user && go test ./...
	cd services/subscription && go test ./...
	cd services/notification && go test ./...
