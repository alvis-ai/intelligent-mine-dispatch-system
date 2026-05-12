package handler

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHub struct {
	rdb     *redis.Client
	clients map[*websocket.Conn]bool
	mu      sync.RWMutex
}

func NewWSHub(rdb *redis.Client) *WSHub {
	return &WSHub{
		rdb:     rdb,
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *WSHub) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		conn.Close()
	}()

	pubsub := h.rdb.Subscribe(r.Context(), "vehicle:location:updates")
	defer pubsub.Close()
	ch := pubsub.Channel()

	for msg := range ch {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			break
		}
	}
}

func (h *WSHub) Broadcast(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		client.WriteMessage(websocket.TextMessage, data)
	}
}

func RegisterWSRoute(server *rest.Server, hub *WSHub) {
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/ws/location",
		Handler: hub.HandleWS,
	})
}

func (h *WSHub) StartRedisListener(ctx context.Context) {
	pubsub := h.rdb.Subscribe(ctx, "vehicle:location:updates")
	defer pubsub.Close()
	ch := pubsub.Channel()
	for msg := range ch {
		h.Broadcast([]byte(msg.Payload))
	}
}
