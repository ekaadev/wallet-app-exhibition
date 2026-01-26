package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// MutationType represents the type of wallet mutation
type MutationType string

const (
	MutationTypeDebit  MutationType = "debit"
	MutationTypeCredit MutationType = "credit"
)

type WalletMutation struct {
	ID            uint            `gorm:"column:id;primaryKey;autoIncrement"`
	WalletID      uint            `gorm:"column:wallet_id;not null"`
	TransactionID uint            `gorm:"column:transaction_id;not null"`
	Type          MutationType    `gorm:"column:type;type:enum('debit','credit');not null"`
	Amount        decimal.Decimal `gorm:"column:amount;type:decimal(20,2);not null"`
	BalanceBefore decimal.Decimal `gorm:"column:balance_before;type:decimal(20,2);not null"`
	BalanceAfter  decimal.Decimal `gorm:"column:balance_after;type:decimal(20,2);not null"`
	CreatedAt     time.Time       `gorm:"column:created_at;autoCreateTime;not null"`

	// Relations
	Wallet      *Wallet      `gorm:"foreignKey:WalletID;references:ID"`
	Transaction *Transaction `gorm:"foreignKey:TransactionID;references:ID"`
}

func (wm *WalletMutation) TableName() string {
	return "wallet_mutations"
}
