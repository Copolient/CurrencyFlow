package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

// MessageType is carried in every WebSocket payload so clients can distinguish
// public rate updates from user-specific notifications without a second socket.
type MessageType string

const (
	MessageTypeRate         MessageType = "rate"
	MessageTypeNotification MessageType = "notification"
)

type RateUpdate struct {
	Type         string  `json:"type"`
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Rate         float64 `json:"rate"`
	Timestamp    string  `json:"timestamp"`
}

type NotificationMessage struct {
	Type             string    `json:"type"`
	ID               uint      `json:"id"`
	NotificationType string    `json:"notificationType"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	Read             bool      `json:"read"`
	CreatedAt        time.Time `json:"createdAt"`
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	pair   string // subscribed currency pair, empty = all
	userID uint   // authenticated user; 0 means anonymous (rate updates only)
	mu     sync.Mutex
	closed bool
}

type Hub struct {
	clients       map[*Client]bool
	broadcast     chan []byte
	userBroadcast chan userMessage
	register      chan *Client
	unregister    chan *Client
	done          chan struct{}
	mu            sync.RWMutex
}

type userMessage struct {
	userID uint
	data   []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		broadcast:     make(chan []byte, 256),
		userBroadcast: make(chan userMessage, 256),
		register:      make(chan *Client),
		unregister:    make(chan *Client, 64),
		done:          make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.done:
			h.mu.Lock()
			for client := range h.clients {
				client.Close()
				delete(h.clients, client)
			}
			h.mu.Unlock()
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			count := len(h.clients)
			h.mu.Unlock()
			log.Printf("WebSocket client connected (total: %d)", count)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			count := len(h.clients)
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected (total: %d)", count)

		case message := <-h.broadcast:
			h.dispatch(message)

		case message := <-h.userBroadcast:
			h.dispatchToUser(message.userID, message.data)
		}
	}
}

// Stop signals the hub to shut down gracefully.
func (h *Hub) Stop() {
	close(h.done)
}

func (h *Hub) BroadcastRateUpdate(update RateUpdate) {
	if update.Type == "" {
		update.Type = string(MessageTypeRate)
	}
	data, err := json.Marshal(update)
	if err != nil {
		log.Printf("Hub: failed to marshal rate update: %v", err)
		return
	}
	h.broadcast <- data
}

// BroadcastNotification sends a user-specific notification to every active
// connection authenticated as userID.
func (h *Hub) BroadcastNotification(userID uint, notification NotificationMessage) {
	notification.Type = string(MessageTypeNotification)
	data, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Hub: failed to marshal notification: %v", err)
		return
	}
	h.userBroadcast <- userMessage{userID: userID, data: data}
}

func (h *Hub) dispatch(message []byte) {
	var stale []*Client
	h.mu.RLock()
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			stale = append(stale, client)
		}
	}
	h.mu.RUnlock()
	for _, c := range stale {
		h.enqueueUnregister(c)
	}
}

func (h *Hub) dispatchToUser(userID uint, message []byte) {
	var stale []*Client
	h.mu.RLock()
	for client := range h.clients {
		if client.userID != userID {
			continue
		}
		select {
		case client.send <- message:
		default:
			stale = append(stale, client)
		}
	}
	h.mu.RUnlock()
	for _, c := range stale {
		h.enqueueUnregister(c)
	}
}

func (h *Hub) enqueueUnregister(client *Client) {
	select {
	case h.unregister <- client:
	default:
		// If the channel is full, ReadPump will retry on connection close;
		// dropping here prevents the hub from blocking during shutdown.
	}
}

func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}
func NewClient(hub *Hub, conn *websocket.Conn, pair string, userID uint) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		pair:   pair,
		userID: userID,
	}
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.closed = true
		close(c.send)
		c.conn.Close()
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.enqueueUnregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Pair filtering must not swallow user notifications.
			if c.pair != "" {
				var envelope struct {
					Type         string `json:"type"`
					FromCurrency string `json:"fromCurrency"`
					ToCurrency   string `json:"toCurrency"`
				}
				if err := json.Unmarshal(message, &envelope); err == nil && envelope.Type == string(MessageTypeRate) {
					pair := envelope.FromCurrency + "/" + envelope.ToCurrency
					if pair != c.pair {
						continue
					}
				}
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
