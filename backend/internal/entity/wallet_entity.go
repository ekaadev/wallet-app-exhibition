package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        uint            `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint            `gorm:"column:user_id;uniqueIndex;not null"`
	Balance   decimal.Decimal `gorm:"column:balance;type:decimal(20,2);not null;default:0.00"`
	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt time.Time       `gorm:"column:updated_at;autoUpdateTime;not null"`

	// Relations
	User *User `gorm:"foreignKey:UserID;references:ID"`
}

func (w *Wallet) TableName() string {
	return "wallets"
}
