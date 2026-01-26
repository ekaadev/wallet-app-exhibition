package websocket

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active WebSocket connections and broadcasts messages to users.
type Hub struct {
	// connections maps user ID to their WebSocket connections
	connections map[uint][]*websocket.Conn
	// mutex for thread-safe operations
	mu sync.RWMutex
	// logger
	log *logrus.Logger
}

// NewHub creates a new WebSocket Hub.
func NewHub(log *logrus.Logger) *Hub {
	return &Hub{
		connections: make(map[uint][]*websocket.Conn),
		log:         log,
	}
}

// Register adds a new connection for a user.
func (h *Hub) Register(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.connections[userID] = append(h.connections[userID], conn)
	h.log.Infof("WebSocket connection registered for user ID: %d, total connections: %d", userID, len(h.connections[userID]))
}

// Unregister removes a connection for a user.
func (h *Hub) Unregister(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns := h.connections[userID]
	for i, c := range conns {
		if c == conn {
			// Remove connection from slice
			h.connections[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	// Clean up empty user entries
	if len(h.connections[userID]) == 0 {
		delete(h.connections, userID)
	}

	h.log.Infof("WebSocket connection unregistered for user ID: %d", userID)
}

// BroadcastToUser sends a message to all connections of a specific user.
func (h *Hub) BroadcastToUser(userID uint, message []byte) error {
	h.mu.RLock()
	conns := h.connections[userID]
	h.mu.RUnlock()

	if len(conns) == 0 {
		h.log.Debugf("No active connections for user ID: %d", userID)
		return nil
	}

	var lastErr error
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			h.log.Warnf("Failed to send message to user ID %d: %v", userID, err)
			lastErr = err
		}
	}

	return lastErr
}

// GetConnectionCount returns the number of active connections for a user.
func (h *Hub) GetConnectionCount(userID uint) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.connections[userID])
}

// GetTotalConnections returns the total number of active connections.
func (h *Hub) GetTotalConnections() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	total := 0
	for _, conns := range h.connections {
		total += len(conns)
	}
	return total
}
