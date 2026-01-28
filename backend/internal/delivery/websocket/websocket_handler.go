package websocket

import (
	"backend/internal/util"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Handler handles WebSocket connections.
type Handler struct {
	Hub       *Hub
	TokenUtil *util.TokenUtil
	Log       *logrus.Logger
}

// NewHandler creates a new WebSocket Handler.
func NewHandler(hub *Hub, tokenUtil *util.TokenUtil, log *logrus.Logger) *Handler {
	return &Handler{
		Hub:       hub,
		TokenUtil: tokenUtil,
		Log:       log,
	}
}

// UpgradeMiddleware checks if the request is a WebSocket upgrade and validates the token.
func (h *Handler) UpgradeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if this is a WebSocket upgrade request
		if websocket.IsWebSocketUpgrade(c) {
			// Get token from query parameter or cookie
			token := c.Query("token")
			if token == "" {
				token = c.Cookies("jwt")
			}

			if token == "" {
				h.Log.Warn("WebSocket connection attempt without token")
				return fiber.ErrUnauthorized
			}

			// Remove Bearer prefix if present
			token = strings.TrimPrefix(token, "Bearer ")

			// Parse and verify the token
			auth, err := h.TokenUtil.ParseToken(c.UserContext(), token)
			if err != nil {
				h.Log.Warnf("WebSocket token validation failed: %v", err)
				return fiber.ErrUnauthorized
			}

			// Store auth in locals for the WebSocket handler
			c.Locals("auth", auth)
			c.Locals("userID", *auth.UserID)

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

// HandleConnection handles WebSocket connections.
func (h *Handler) HandleConnection() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// Get user ID from locals
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			h.Log.Warn("Failed to get user ID from locals")
			return
		}

		// Register connection
		h.Hub.Register(userID, c)
		defer h.Hub.Unregister(userID, c)

		h.Log.Infof("WebSocket connection established for user ID: %d", userID)

		// Keep connection alive and handle messages
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					h.Log.Warnf("WebSocket unexpected close error for user ID %d: %v", userID, err)
				}
				break
			}

			// Handle ping/pong for heartbeat
			if messageType == websocket.PingMessage {
				if err := c.WriteMessage(websocket.PongMessage, nil); err != nil {
					h.Log.Warnf("Failed to send pong to user ID %d: %v", userID, err)
					break
				}
				continue
			}

			// Log received messages (for debugging)
			h.Log.Debugf("Received message from user ID %d: %s", userID, string(message))
		}

		h.Log.Infof("WebSocket connection closed for user ID: %d", userID)
	})
}
