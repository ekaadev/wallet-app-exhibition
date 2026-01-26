package websocket

import (
	"backend/internal/model"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

// NotifierInterface defines the interface for notification operations.
type NotifierInterface interface {
	NotifyTransaction(userID uint, notification *model.TransactionNotification) error
	NotifyWalletUpdate(userID uint, notification *model.WalletUpdateNotification) error
}

// Notifier sends notifications to users via WebSocket.
type Notifier struct {
	Hub *Hub
	Log *logrus.Logger
}

// NewNotifier creates a new Notifier instance.
func NewNotifier(hub *Hub, log *logrus.Logger) *Notifier {
	return &Notifier{
		Hub: hub,
		Log: log,
	}
}

// NotifyTransaction sends a transaction notification to a user.
func (n *Notifier) NotifyTransaction(userID uint, notification *model.TransactionNotification) error {
	message := model.WebSocketMessage{
		Type:    "transaction",
		Payload: notification,
	}

	data, err := json.Marshal(message)
	if err != nil {
		n.Log.Errorf("Failed to marshal transaction notification: %v", err)
		return err
	}

	if err := n.Hub.BroadcastToUser(userID, data); err != nil {
		n.Log.Warnf("Failed to send transaction notification to user ID %d: %v", userID, err)
		return err
	}

	n.Log.Infof("Transaction notification sent to user ID: %d", userID)
	return nil
}

// NotifyWalletUpdate sends a wallet update notification to a user.
func (n *Notifier) NotifyWalletUpdate(userID uint, notification *model.WalletUpdateNotification) error {
	if notification.UpdatedAt.IsZero() {
		notification.UpdatedAt = time.Now()
	}

	message := model.WebSocketMessage{
		Type:    "wallet_update",
		Payload: notification,
	}

	data, err := json.Marshal(message)
	if err != nil {
		n.Log.Errorf("Failed to marshal wallet update notification: %v", err)
		return err
	}

	if err := n.Hub.BroadcastToUser(userID, data); err != nil {
		n.Log.Warnf("Failed to send wallet update notification to user ID %d: %v", userID, err)
		return err
	}

	n.Log.Infof("Wallet update notification sent to user ID: %d", userID)
	return nil
}
