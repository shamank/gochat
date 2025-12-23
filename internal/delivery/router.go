package delivery

import (
	"net/http"

	"gochat/internal/delivery/handler"
	"gochat/internal/delivery/websocket"
)

type Router struct {
	userHandler    *handler.UserHandler
	roomHandler    *handler.RoomHandler
	messageHandler *handler.MessageHandler
	wsHub          *websocket.Hub
}

func NewRouter(
	userHandler *handler.UserHandler,
	roomHandler *handler.RoomHandler,
	messageHandler *handler.MessageHandler,
	wsHub *websocket.Hub,
) *Router {
	return &Router{
		userHandler:    userHandler,
		roomHandler:    roomHandler,
		messageHandler: messageHandler,
		wsHub:          wsHub,
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/users/register", r.userHandler.RegisterUser)
	mux.HandleFunc("/api/users/get", r.userHandler.GetUser)

	mux.HandleFunc("/api/rooms/create", r.roomHandler.CreateRoom)
	mux.HandleFunc("/api/rooms/get", r.roomHandler.GetRoom)
	mux.HandleFunc("/api/rooms/all", r.roomHandler.GetAllRooms)

	mux.HandleFunc("/api/messages/send", r.messageHandler.SendMessage)
	mux.HandleFunc("/api/messages/history", r.messageHandler.GetMessagesHistory)

	wsHub := r.wsHub
	mux.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		websocket.ServeWS(wsHub, w, req)
	})

	return mux
}
