package websocket_test

import (
	"backend/internal/delivery/websocket"
	"backend/internal/model"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// createTestHub creates a Hub instance for testing.
func createTestHub() *websocket.Hub {
	log := logrus.New()
	log.SetOutput(io.Discard)
	return websocket.NewHub(log)
}

// TestHub_NewHub tests creating a new Hub instance.
func TestHub_NewHub(t *testing.T) {
	hub := createTestHub()
	assert.NotNil(t, hub)
	assert.Equal(t, 0, hub.GetTotalConnections())
}

// TestHub_GetConnectionCount_NoConnections tests getting connection count when no connections exist.
func TestHub_GetConnectionCount_NoConnections(t *testing.T) {
	hub := createTestHub()

	count := hub.GetConnectionCount(1)
	assert.Equal(t, 0, count)
}

// TestHub_GetTotalConnections_NoConnections tests getting total connections when no connections exist.
func TestHub_GetTotalConnections_NoConnections(t *testing.T) {
	hub := createTestHub()

	total := hub.GetTotalConnections()
	assert.Equal(t, 0, total)
}

// TestHub_BroadcastToNonExistentUser tests broadcasting to a user with no connections.
func TestHub_BroadcastToNonExistentUser(t *testing.T) {
	hub := createTestHub()

	err := hub.BroadcastToUser(999, []byte(`{"type":"test"}`))
	// Should not return error, just log that no connections exist
	assert.NoError(t, err)
}

// TestNotifier_NewNotifier tests creating a new Notifier instance.
func TestNotifier_NewNotifier(t *testing.T) {
	hub := createTestHub()
	log := logrus.New()
	log.SetOutput(io.Discard)

	notifier := websocket.NewNotifier(hub, log)
	assert.NotNil(t, notifier)
}

// TestNotifier_NotifyTransaction_NoConnections tests notifying transaction when user has no connections.
func TestNotifier_NotifyTransaction_NoConnections(t *testing.T) {
	hub := createTestHub()
	log := logrus.New()
	log.SetOutput(io.Discard)

	notifier := websocket.NewNotifier(hub, log)

	notification := &model.TransactionNotification{
		TransactionID:     1,
		TransactionType:   "transfer",
		Amount:            "50000",
		ToUserID:          2,
		PerformedByUserID: 1,
		CreatedAt:         "2026-01-26T12:00:00Z",
	}

	// Should not return error even when no connections exist
	err := notifier.NotifyTransaction(999, notification)
	assert.NoError(t, err)
}

// TestNotifier_NotifyWalletUpdate_NoConnections tests notifying wallet update when user has no connections.
func TestNotifier_NotifyWalletUpdate_NoConnections(t *testing.T) {
	hub := createTestHub()
	log := logrus.New()
	log.SetOutput(io.Discard)

	notifier := websocket.NewNotifier(hub, log)

	notification := &model.WalletUpdateNotification{
		WalletID:      1,
		NewBalance:    "150000",
		MutationType:  "credit",
		MutationID:    1,
		TransactionID: 1,
		Amount:        "50000",
	}

	// Should not return error even when no connections exist
	err := notifier.NotifyWalletUpdate(999, notification)
	assert.NoError(t, err)
}
