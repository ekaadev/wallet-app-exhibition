package usecase

import (
	"backend/internal/model"
	"backend/internal/model/converter"
	"backend/internal/repository"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletMutationUseCase struct {
	DB                       *gorm.DB
	Log                      *logrus.Logger
	Validate                 *validator.Validate
	WalletMutationRepository *repository.WalletMutationRepository
	WalletRepository         *repository.WalletRepository
}

func NewWalletMutationUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	walletMutationRepo *repository.WalletMutationRepository,
	walletRepo *repository.WalletRepository,
) *WalletMutationUseCase {
	return &WalletMutationUseCase{
		DB:                       db,
		Log:                      log,
		Validate:                 validate,
		WalletMutationRepository: walletMutationRepo,
		WalletRepository:         walletRepo,
	}
}

// GetMutationsByUserID retrieves wallet mutations for a user.
func (uc *WalletMutationUseCase) GetMutationsByUserID(ctx context.Context, userID uint, page, limit int) (*model.WalletMutationListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Get wallet for user
	wallet, err := uc.WalletRepository.FindByUserID(uc.DB.WithContext(ctx), userID)
	if err != nil {
		uc.Log.Errorf("FindByUserID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if wallet == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Wallet not found")
	}

	mutations, total, err := uc.WalletMutationRepository.FindByWalletID(uc.DB.WithContext(ctx), wallet.ID, page, limit)
	if err != nil {
		uc.Log.Errorf("FindByWalletID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.WalletMutationListResponse{
		Mutations: converter.WalletMutationsToWalletMutationResponses(mutations),
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}
