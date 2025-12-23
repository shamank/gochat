.PHONY: server client test

SERVER_CMD = cmd/server
CLIENT_CMD = client
PORT ?= 8080
HOST ?= 0.0.0.0
SERVER_URL ?= http://localhost:8080

server: ## Запустить сервер
	@set PORT=$(PORT) && set HOST=$(HOST) && go run $(SERVER_CMD)/main.go

client: ## Запустить клиент
	@set SERVER_URL=$(SERVER_URL) && go run ./$(CLIENT_CMD)

test: ## Запустить тесты
	@go test -v ./internal/...