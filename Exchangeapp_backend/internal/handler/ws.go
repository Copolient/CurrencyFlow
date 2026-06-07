package handler

import (
	"log"
	"net/http"
	"os"
	"strings"

	ws "exchangeapp/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func newUpgrader() websocket.Upgrader {
	allowed := os.Getenv("WS_ALLOWED_ORIGINS")
	if allowed == "" {
		allowed = "http://localhost:5173,http://localhost:80"
	}
	origins := make(map[string]bool)
	for _, o := range strings.Split(allowed, ",") {
		origins[strings.TrimSpace(o)] = true
	}

	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origins[origin]
		},
	}
}

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	upgrader := newUpgrader()
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	pair := c.Query("pair") // optional: subscribe to specific currency pair

	client := ws.NewClient(h.hub, conn, pair)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func (h *WSHandler) GetClientCount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"clients": h.hub.ClientCount()})
}
