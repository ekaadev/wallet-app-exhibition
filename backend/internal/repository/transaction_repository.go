package repository

import (
	"backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	Repository[entity.Transaction]
	Log *logrus.Logger
}

func NewTransactionRepository(log *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{
		Log: log,
	}
}

func (r *TransactionRepository) FindByWalletID(db *gorm.DB, walletID uint, page, limit int) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var total int64

	query := db.Model(&entity.Transaction{}).
		Where("from_wallet_id = ? OR to_wallet_id = ?", walletID, walletID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepository) FindByUserID(db *gorm.DB, userID uint, page, limit int) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var total int64

	query := db.Model(&entity.Transaction{}).
		Where("performed_by_user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}
