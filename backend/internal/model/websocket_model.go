package model

import "time"

// WebSocketMessage represents a generic WebSocket message.
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// TransactionNotification represents a notification for transaction events.
type TransactionNotification struct {
	TransactionID     uint    `json:"transaction_id"`
	TransactionType   string  `json:"transaction_type"`
	Amount            string  `json:"amount"`
	FromUserID        *uint   `json:"from_user_id,omitempty"`
	ToUserID          uint    `json:"to_user_id"`
	PerformedByUserID uint    `json:"performed_by_user_id"`
	Description       *string `json:"description,omitempty"`
	CreatedAt         string  `json:"created_at"`
}

// WalletUpdateNotification represents a notification for wallet balance updates.
type WalletUpdateNotification struct {
	WalletID      uint      `json:"wallet_id"`
	NewBalance    string    `json:"new_balance"`
	MutationType  string    `json:"mutation_type"`
	MutationID    uint      `json:"mutation_id"`
	TransactionID uint      `json:"transaction_id"`
	Amount        string    `json:"amount"`
	UpdatedAt     time.Time `json:"updated_at"`
}
