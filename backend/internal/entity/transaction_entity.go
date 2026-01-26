package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeTopUp    TransactionType = "top_up"
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypeWithdraw TransactionType = "withdraw"
)

// TransactionStatus represents the status of transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID                uint              `gorm:"column:id;primaryKey;autoIncrement"`
	Type              TransactionType   `gorm:"column:type;type:enum('top_up','transfer','withdraw');not null"`
	Amount            decimal.Decimal   `gorm:"column:amount;type:decimal(20,2);not null"`
	FromWalletID      *uint             `gorm:"column:from_wallet_id"`
	ToWalletID        uint              `gorm:"column:to_wallet_id;not null"`
	PerformedByUserID uint              `gorm:"column:performed_by_user_id;not null"`
	Status            TransactionStatus `gorm:"column:status;type:enum('pending','completed','failed');not null;default:'pending'"`
	Description       *string           `gorm:"column:description;type:varchar(255)"`
	CreatedAt         time.Time         `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt         time.Time         `gorm:"column:updated_at;autoUpdateTime;not null"`

	// Relations
	FromWallet      *Wallet `gorm:"foreignKey:FromWalletID;references:ID"`
	ToWallet        *Wallet `gorm:"foreignKey:ToWalletID;references:ID"`
	PerformedByUser *User   `gorm:"foreignKey:PerformedByUserID;references:ID"`
}

func (t *Transaction) TableName() string {
	return "transactions"
}
