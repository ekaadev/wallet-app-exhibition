package repository

import (
	"backend/internal/entity"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletRepository struct {
	Repository[entity.Wallet]
	Log *logrus.Logger
}

func NewWalletRepository(log *logrus.Logger) *WalletRepository {
	return &WalletRepository{
		Log: log,
	}
}

func (r *WalletRepository) FindByUserID(db *gorm.DB, userID uint) (*entity.Wallet, error) {
	var wallet entity.Wallet
	err := db.Where("user_id = ?", userID).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &wallet, err
}

func (r *WalletRepository) UpdateBalance(db *gorm.DB, walletID uint, newBalance interface{}) error {
	return db.Model(&entity.Wallet{}).Where("id = ?", walletID).Update("balance", newBalance).Error
}

func (r *WalletRepository) LockForUpdate(db *gorm.DB, walletID uint) (*entity.Wallet, error) {
	var wallet entity.Wallet
	err := db.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", walletID).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &wallet, err
}
