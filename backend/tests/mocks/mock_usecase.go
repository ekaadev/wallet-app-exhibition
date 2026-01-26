package mocks

import (
	"backend/internal/entity"
	"backend/internal/model"
	"context"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserUseCase is a mock implementation of UserUseCaseInterface.
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Create(ctx context.Context, request *model.UserRegistrationRequest) (*model.UserResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) Login(ctx context.Context, request *model.UserLoginRequest) (*model.UserResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) GetProfile(ctx context.Context, userID uint) (*model.UserProfileResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserProfileResponse), args.Error(1)
}

// MockWalletUseCase is a mock implementation of WalletUseCaseInterface.
type MockWalletUseCase struct {
	mock.Mock
}

func (m *MockWalletUseCase) GetByUserID(ctx context.Context, userID uint) (*model.WalletResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WalletResponse), args.Error(1)
}

func (m *MockWalletUseCase) Create(ctx context.Context, tx *gorm.DB, userID uint) (*entity.Wallet, error) {
	args := m.Called(ctx, tx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Wallet), args.Error(1)
}

// MockTransactionUseCase is a mock implementation of TransactionUseCaseInterface.
type MockTransactionUseCase struct {
	mock.Mock
}

func (m *MockTransactionUseCase) TopUp(ctx context.Context, auth *model.Auth, request *model.TopUpRequest) (*model.TransactionResponse, error) {
	args := m.Called(ctx, auth, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TransactionResponse), args.Error(1)
}

func (m *MockTransactionUseCase) Transfer(ctx context.Context, auth *model.Auth, request *model.TransferRequest) (*model.TransactionResponse, error) {
	args := m.Called(ctx, auth, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TransactionResponse), args.Error(1)
}

func (m *MockTransactionUseCase) GetTransactionsByUserID(ctx context.Context, userID uint, page, limit int) (*model.TransactionListResponse, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TransactionListResponse), args.Error(1)
}

// MockWalletMutationUseCase is a mock implementation of WalletMutationUseCaseInterface.
type MockWalletMutationUseCase struct {
	mock.Mock
}

func (m *MockWalletMutationUseCase) GetMutationsByUserID(ctx context.Context, userID uint, page, limit int) (*model.WalletMutationListResponse, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WalletMutationListResponse), args.Error(1)
}
