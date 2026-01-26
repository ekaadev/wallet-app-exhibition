package usecase

import (
	"backend/internal/entity"
	"backend/internal/model"
	"context"

	"gorm.io/gorm"
)

// UserUseCaseInterface defines the interface for user-related use cases.
type UserUseCaseInterface interface {
	Create(ctx context.Context, request *model.UserRegistrationRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.UserLoginRequest) (*model.UserResponse, error)
	GetProfile(ctx context.Context, userID uint) (*model.UserProfileResponse, error)
}

// WalletUseCaseInterface defines the interface for wallet-related use cases.
type WalletUseCaseInterface interface {
	GetByUserID(ctx context.Context, userID uint) (*model.WalletResponse, error)
	Create(ctx context.Context, tx *gorm.DB, userID uint) (*entity.Wallet, error)
}

// TransactionUseCaseInterface defines the interface for transaction-related use cases.
type TransactionUseCaseInterface interface {
	TopUp(ctx context.Context, auth *model.Auth, request *model.TopUpRequest) (*model.TransactionResponse, error)
	Transfer(ctx context.Context, auth *model.Auth, request *model.TransferRequest) (*model.TransactionResponse, error)
	GetTransactionsByUserID(ctx context.Context, userID uint, page, limit int) (*model.TransactionListResponse, error)
}

// WalletMutationUseCaseInterface defines the interface for wallet mutation-related use cases.
type WalletMutationUseCaseInterface interface {
	GetMutationsByUserID(ctx context.Context, userID uint, page, limit int) (*model.WalletMutationListResponse, error)
}
