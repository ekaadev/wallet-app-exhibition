package usecase

import (
	"backend/internal/entity"
	"backend/internal/model"
	"backend/internal/model/converter"
	"backend/internal/repository"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	WalletRepository *repository.WalletRepository
}

func NewWalletUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, walletRepo *repository.WalletRepository) *WalletUseCase {
	return &WalletUseCase{
		DB:               db,
		Log:              log,
		Validate:         validate,
		WalletRepository: walletRepo,
	}
}

// GetByUserID retrieves the wallet for a specific user.
func (uc *WalletUseCase) GetByUserID(ctx context.Context, userID uint) (*model.WalletResponse, error) {
	wallet, err := uc.WalletRepository.FindByUserID(uc.DB.WithContext(ctx), userID)
	if err != nil {
		uc.Log.Errorf("FindByUserID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if wallet == nil {
		uc.Log.Warnf("Wallet not found for user ID: %d", userID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Wallet not found")
	}

	return converter.WalletToWalletResponse(wallet), nil
}

// Create creates a new wallet for a user (used during registration).
func (uc *WalletUseCase) Create(ctx context.Context, tx *gorm.DB, userID uint) (*entity.Wallet, error) {
	wallet := &entity.Wallet{
		UserID: userID,
	}

	if err := uc.WalletRepository.Create(tx, wallet); err != nil {
		uc.Log.Errorf("Wallet creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return wallet, nil
}
