package handler

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"currencyflow/internal/repository"
	ws "currencyflow/internal/websocket"
	"currencyflow/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	defaultUpgrader *websocket.Upgrader
	upgraderOnce    sync.Once
)

func getDefaultUpgrader() *websocket.Upgrader {
	upgraderOnce.Do(func() {
		allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
		if allowed == "" {
			allowed = os.Getenv("WS_ALLOWED_ORIGINS")
		}
		if allowed == "" {
			allowed = "http://localhost:5173,http://localhost,http://localhost:80"
		}
		origins := make(map[string]bool)
		for _, o := range strings.Split(allowed, ",") {
			origins[strings.TrimSpace(o)] = true
		}

		defaultUpgrader = &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return origins[origin]
			},
		}
	})
	return defaultUpgrader
}

type WSHandler struct {
	hub      *ws.Hub
	jwt      *auth.JWTManager
	userRepo repository.UserRepository
}

func NewWSHandler(hub *ws.Hub, jwt *auth.JWTManager, userRepo repository.UserRepository) *WSHandler {
	return &WSHandler{hub: hub, jwt: jwt, userRepo: userRepo}
}

// resolveUserID returns the authenticated user's ID for a token passed as a
// query parameter. Browsers cannot set custom headers on WebSocket handshakes,
// so the token is passed via ?token=. Anonymous connections are allowed and
// still receive public rate updates.
func (h *WSHandler) resolveUserID(c *gin.Context) uint {
	token := c.Query("token")
	if token == "" || h.jwt == nil || h.userRepo == nil {
		return 0
	}
	username, err := h.jwt.ParseToken(token)
	if err != nil || username == "" {
		return 0
	}
	user, err := h.userRepo.FindByUsername(username)
	if err != nil || user == nil {
		return 0
	}
	return user.ID
}

func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	// Limit concurrent WebSocket connections
	if h.hub.ClientCount() >= 1000 {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "too many connections"})
		return
	}

	conn, err := getDefaultUpgrader().Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	pair := c.Query("pair") // optional: subscribe to specific currency pair
	userID := h.resolveUserID(c)

	client := ws.NewClient(h.hub, conn, pair, userID)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func (h *WSHandler) GetClientCount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"clients": h.hub.ClientCount()})
}
