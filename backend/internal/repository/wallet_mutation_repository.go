package repository

import (
	"backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletMutationRepository struct {
	Repository[entity.WalletMutation]
	Log *logrus.Logger
}

func NewWalletMutationRepository(log *logrus.Logger) *WalletMutationRepository {
	return &WalletMutationRepository{
		Log: log,
	}
}

func (r *WalletMutationRepository) FindByWalletID(db *gorm.DB, walletID uint, page, limit int) ([]entity.WalletMutation, int64, error) {
	var mutations []entity.WalletMutation
	var total int64

	query := db.Model(&entity.WalletMutation{}).Where("wallet_id = ?", walletID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&mutations).Error
	if err != nil {
		return nil, 0, err
	}

	return mutations, total, nil
}

func (r *WalletMutationRepository) FindByTransactionID(db *gorm.DB, transactionID uint) ([]entity.WalletMutation, error) {
	var mutations []entity.WalletMutation
	err := db.Where("transaction_id = ?", transactionID).Find(&mutations).Error
	return mutations, err
}
