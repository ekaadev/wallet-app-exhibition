package usecase

import (
	"backend/internal/entity"
	"backend/internal/model"
	"backend/internal/model/converter"
	"backend/internal/repository"
	"backend/internal/util"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	UserRepository   *repository.UserRepository
	WalletRepository *repository.WalletRepository
	TokenUtil        *util.TokenUtil
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepo *repository.UserRepository, walletRepo *repository.WalletRepository, tokenUtil *util.TokenUtil) *UserUseCase {
	return &UserUseCase{
		DB:               db,
		Log:              log,
		Validate:         validate,
		UserRepository:   userRepo,
		WalletRepository: walletRepo,
		TokenUtil:        tokenUtil,
	}
}

// Create handles the business logic for creating a new user.
func (uc *UserUseCase) Create(ctx context.Context, request *model.UserRegistrationRequest) (*model.UserResponse, error) {
	// Start transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	err := uc.Validate.Struct(request)
	if err != nil {
		uc.Log.Warnf("Validation error: %v", err)
		return nil, err
	}

	// Business logic for creating user
	existingUser, err := uc.UserRepository.FindByUsername(tx, request.Username)
	if err != nil {
		uc.Log.Warnf("FindByUsername error: %v", err)
	}

	if existingUser != nil {
		if existingUser.Username == request.Username {
			uc.Log.Warnf("Username already exists")
			return nil, fiber.NewError(fiber.StatusConflict, "Username already exists")
		}
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.Log.Errorf("Password hashing error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		Username: request.Username,
		Password: string(password),
	}

	if err = uc.UserRepository.Create(tx, user); err != nil {
		uc.Log.Warnf("User creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Auto-create wallet for the user with balance 0
	wallet := &entity.Wallet{
		UserID: user.ID,
	}
	if err = uc.WalletRepository.Create(tx, wallet); err != nil {
		uc.Log.Warnf("Wallet creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err = tx.Commit().Error; err != nil {
		uc.Log.Errorf("Transaction commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create token
	token, err := uc.TokenUtil.CreateToken(ctx, &model.Auth{
		UserID:   &user.ID,
		Username: request.Username,
		Role:     user.Role,
	})
	if err != nil {
		uc.Log.Warnf("Token creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToUserResponse(user, token), err
}

// Login handles the business logic for user login.
func (uc *UserUseCase) Login(ctx context.Context, request *model.UserLoginRequest) (*model.UserResponse, error) {
	// Validate request
	err := uc.Validate.Struct(request)
	if err != nil {
		uc.Log.Warnf("Validation error: %v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Find user by username
	user, err := uc.UserRepository.FindByUsername(uc.DB.WithContext(ctx), request.Username)
	if err != nil {
		uc.Log.Warnf("FindByUsername error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if user == nil {
		uc.Log.Warnf("User not found: %s", request.Username)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		uc.Log.Warnf("Password mismatch for user: %s", request.Username)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	// Create token
	token, err := uc.TokenUtil.CreateToken(ctx, &model.Auth{
		UserID:   &user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		uc.Log.Warnf("Token creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToUserResponse(user, token), nil
}

// GetProfile retrieves user profile with wallet information.
func (uc *UserUseCase) GetProfile(ctx context.Context, userID uint) (*model.UserProfileResponse, error) {
	// Find user by ID
	user, err := uc.UserRepository.FindByID(uc.DB.WithContext(ctx), userID)
	if err != nil {
		uc.Log.Errorf("FindByID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if user == nil {
		uc.Log.Warnf("User not found: %d", userID)
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Find wallet for user
	wallet, err := uc.WalletRepository.FindByUserID(uc.DB.WithContext(ctx), userID)
	if err != nil {
		uc.Log.Errorf("FindByUserID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	response := &model.UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}

	if wallet != nil {
		response.Wallet = &model.UserProfileWalletInfo{
			ID:      wallet.ID,
			Balance: wallet.Balance,
		}
	}

	return response, nil
}
