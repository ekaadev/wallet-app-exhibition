package model

import "github.com/shopspring/decimal"

// WalletResponse represents the response payload for wallet-related operations.
type WalletResponse struct {
	ID      uint            `json:"id"`
	UserID  uint            `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
}

// GetWalletRequest represents the request for getting wallet information.
type GetWalletRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}
