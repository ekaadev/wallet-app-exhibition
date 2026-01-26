package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// TopUpRequest represents the request payload for top-up operation (super admin only).
type TopUpRequest struct {
	ToUserID    uint            `json:"to_user_id" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required"`
	Description string          `json:"description"`
}

// TransferRequest represents the request payload for transfer operation.
type TransferRequest struct {
	ToUserID    uint            `json:"to_user_id" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required"`
	Description string          `json:"description"`
}

// TransactionResponse represents the response payload for transaction-related operations.
type TransactionResponse struct {
	ID                uint            `json:"id"`
	Type              string          `json:"type"`
	Amount            decimal.Decimal `json:"amount"`
	FromWalletID      *uint           `json:"from_wallet_id,omitempty"`
	ToWalletID        uint            `json:"to_wallet_id"`
	PerformedByUserID uint            `json:"performed_by_user_id"`
	Status            string          `json:"status"`
	Description       *string         `json:"description,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
}

// TransactionListRequest represents the request for listing transactions.
type TransactionListRequest struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

// TransactionListResponse represents the response payload for transaction list.
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	Limit        int                   `json:"limit"`
}
