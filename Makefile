.PHONY: server client test

SERVER_CMD = cmd/server
CLIENT_CMD = client

server: ## Запустить сервер
	go run $(SERVER_CMD)/main.go

client: ## Запустить клиент
	go run ./$(CLIENT_CMD)

test: ## Запустить тесты
	@go test -v ./internal/...