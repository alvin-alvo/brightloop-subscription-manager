.PHONY: up down logs tidy build test

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

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
