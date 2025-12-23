package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gochat/internal/delivery"
	"gochat/internal/delivery/handler"
	"gochat/internal/delivery/websocket"
	"gochat/internal/repository"
	"gochat/internal/usecase"
)

func main() {
	userRepo := repository.NewInMemoryUserRepository()
	roomRepo := repository.NewInMemoryRoomRepository()
	messageRepo := repository.NewInMemoryMessageRepository()

	userUsecase := usecase.NewUserUsecase(userRepo)
	roomUsecase := usecase.NewRoomUsecase(roomRepo)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, userRepo, roomRepo)

	wsHub := websocket.NewHub()
	go wsHub.Run()

	userHandler := handler.NewUserHandler(userUsecase)
	roomHandler := handler.NewRoomHandler(roomUsecase)
	messageHandler := handler.NewMessageHandler(messageUsecase, wsHub)

	router := delivery.NewRouter(userHandler, roomHandler, messageHandler, wsHub)
	httpHandler := router.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server starting on %s", addr)
	log.Printf("HTTP API: http://%s", addr)
	log.Printf("WebSocket endpoint: ws://%s/ws", addr)

	if err := http.ListenAndServe(addr, httpHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
