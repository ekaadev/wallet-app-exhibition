package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// WalletMutationResponse represents the response payload for wallet mutation.
type WalletMutationResponse struct {
	ID            uint            `json:"id"`
	WalletID      uint            `json:"wallet_id"`
	TransactionID uint            `json:"transaction_id"`
	Type          string          `json:"type"`
	Amount        decimal.Decimal `json:"amount"`
	BalanceBefore decimal.Decimal `json:"balance_before"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
	CreatedAt     time.Time       `json:"created_at"`
}

// WalletMutationListRequest represents the request for listing wallet mutations.
type WalletMutationListRequest struct {
	WalletID uint `json:"wallet_id" validate:"required"`
	Page     int  `json:"page" validate:"min=1"`
	Limit    int  `json:"limit" validate:"min=1,max=100"`
}

// WalletMutationListResponse represents the response payload for wallet mutation list.
type WalletMutationListResponse struct {
	Mutations []WalletMutationResponse `json:"mutations"`
	Total     int64                    `json:"total"`
	Page      int                      `json:"page"`
	Limit     int                      `json:"limit"`
}
