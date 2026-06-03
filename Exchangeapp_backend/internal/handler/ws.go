package handler

import (
	"log"
	"net/http"

	ws "exchangeapp/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		return true
	},
}

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) HandleWebSocket(c *gin.Context) {
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
