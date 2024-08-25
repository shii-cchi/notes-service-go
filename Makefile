.PHONY: build run
.DEFAULT_GOAL := run

include .env

build:
	go build -o notes_server cmd/main.go

run: build
	./notes_server

migration:
	cd ./internal/database/migrations && goose postgres postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable up

migration_down:
	cd ./internal/database/migrations && goose postgres postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable down

sqlc:
	sqlc generate