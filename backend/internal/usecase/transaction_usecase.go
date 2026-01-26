package usecase

import (
	"backend/internal/delivery/websocket"
	"backend/internal/entity"
	"backend/internal/model"
	"backend/internal/model/converter"
	"backend/internal/repository"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionUseCase struct {
	DB                       *gorm.DB
	Log                      *logrus.Logger
	Validate                 *validator.Validate
	TransactionRepository    *repository.TransactionRepository
	WalletRepository         *repository.WalletRepository
	WalletMutationRepository *repository.WalletMutationRepository
	Notifier                 websocket.NotifierInterface
}

func NewTransactionUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	transactionRepo *repository.TransactionRepository,
	walletRepo *repository.WalletRepository,
	walletMutationRepo *repository.WalletMutationRepository,
) *TransactionUseCase {
	return &TransactionUseCase{
		DB:                       db,
		Log:                      log,
		Validate:                 validate,
		TransactionRepository:    transactionRepo,
		WalletRepository:         walletRepo,
		WalletMutationRepository: walletMutationRepo,
	}
}

// SetNotifier sets the WebSocket notifier for real-time notifications.
func (uc *TransactionUseCase) SetNotifier(notifier websocket.NotifierInterface) {
	uc.Notifier = notifier
}

// TopUp handles top-up operation by super admin.
func (uc *TransactionUseCase) TopUp(ctx context.Context, auth *model.Auth, request *model.TopUpRequest) (*model.TransactionResponse, error) {
	// Validate request
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warnf("Validation error: %v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check if user is super admin
	if auth.Role != "super_admin" {
		uc.Log.Warnf("Unauthorized top-up attempt by user ID: %d", *auth.UserID)
		return nil, fiber.NewError(fiber.StatusForbidden, "Only super admin can perform top-up")
	}

	// Validate amount is positive
	if request.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Amount must be greater than zero")
	}

	// Start transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Find recipient wallet
	toWallet, err := uc.WalletRepository.FindByUserID(tx, request.ToUserID)
	if err != nil {
		uc.Log.Errorf("FindByUserID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if toWallet == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Recipient wallet not found")
	}

	// Lock wallet for update
	toWallet, err = uc.WalletRepository.LockForUpdate(tx, toWallet.ID)
	if err != nil {
		uc.Log.Errorf("LockForUpdate error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create transaction record
	description := request.Description
	transaction := &entity.Transaction{
		Type:              entity.TransactionTypeTopUp,
		Amount:            request.Amount,
		FromWalletID:      nil, // Top-up has no source wallet
		ToWalletID:        toWallet.ID,
		PerformedByUserID: *auth.UserID,
		Status:            entity.TransactionStatusCompleted,
		Description:       &description,
	}

	if err := uc.TransactionRepository.Create(tx, transaction); err != nil {
		uc.Log.Errorf("Transaction creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create wallet mutation for credit
	balanceBefore := toWallet.Balance
	balanceAfter := toWallet.Balance.Add(request.Amount)

	mutation := &entity.WalletMutation{
		WalletID:      toWallet.ID,
		TransactionID: transaction.ID,
		Type:          entity.MutationTypeCredit,
		Amount:        request.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
	}

	if err := uc.WalletMutationRepository.Create(tx, mutation); err != nil {
		uc.Log.Errorf("Wallet mutation creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Update wallet balance
	if err := uc.WalletRepository.UpdateBalance(tx, toWallet.ID, balanceAfter); err != nil {
		uc.Log.Errorf("UpdateBalance error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Errorf("Transaction commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Send real-time notification to recipient
	if uc.Notifier != nil {
		go func() {
			notification := &model.TransactionNotification{
				TransactionID:     transaction.ID,
				TransactionType:   string(transaction.Type),
				Amount:            request.Amount.String(),
				ToUserID:          request.ToUserID,
				PerformedByUserID: *auth.UserID,
				Description:       transaction.Description,
				CreatedAt:         time.Now().Format(time.RFC3339),
			}
			uc.Notifier.NotifyTransaction(request.ToUserID, notification)

			// Also notify wallet update
			walletNotification := &model.WalletUpdateNotification{
				WalletID:      toWallet.ID,
				NewBalance:    balanceAfter.String(),
				MutationType:  string(entity.MutationTypeCredit),
				MutationID:    mutation.ID,
				TransactionID: transaction.ID,
				Amount:        request.Amount.String(),
			}
			uc.Notifier.NotifyWalletUpdate(request.ToUserID, walletNotification)
		}()
	}

	return converter.TransactionToTransactionResponse(transaction), nil
}

// Transfer handles transfer between users.
func (uc *TransactionUseCase) Transfer(ctx context.Context, auth *model.Auth, request *model.TransferRequest) (*model.TransactionResponse, error) {
	// Validate request
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warnf("Validation error: %v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate amount is positive
	if request.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Amount must be greater than zero")
	}

	// Cannot transfer to self
	if *auth.UserID == request.ToUserID {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Cannot transfer to yourself")
	}

	// Start transaction
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Find sender wallet
	fromWallet, err := uc.WalletRepository.FindByUserID(tx, *auth.UserID)
	if err != nil {
		uc.Log.Errorf("FindByUserID error for sender: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if fromWallet == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Your wallet not found")
	}

	// Find recipient wallet
	toWallet, err := uc.WalletRepository.FindByUserID(tx, request.ToUserID)
	if err != nil {
		uc.Log.Errorf("FindByUserID error for recipient: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if toWallet == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Recipient wallet not found")
	}

	// Lock wallets for update (lock in consistent order to prevent deadlock)
	var firstWallet, secondWallet *entity.Wallet
	if fromWallet.ID < toWallet.ID {
		firstWallet, err = uc.WalletRepository.LockForUpdate(tx, fromWallet.ID)
		if err != nil {
			uc.Log.Errorf("LockForUpdate error: %v", err)
			return nil, fiber.ErrInternalServerError
		}
		secondWallet, err = uc.WalletRepository.LockForUpdate(tx, toWallet.ID)
		if err != nil {
			uc.Log.Errorf("LockForUpdate error: %v", err)
			return nil, fiber.ErrInternalServerError
		}
		fromWallet = firstWallet
		toWallet = secondWallet
	} else {
		firstWallet, err = uc.WalletRepository.LockForUpdate(tx, toWallet.ID)
		if err != nil {
			uc.Log.Errorf("LockForUpdate error: %v", err)
			return nil, fiber.ErrInternalServerError
		}
		secondWallet, err = uc.WalletRepository.LockForUpdate(tx, fromWallet.ID)
		if err != nil {
			uc.Log.Errorf("LockForUpdate error: %v", err)
			return nil, fiber.ErrInternalServerError
		}
		toWallet = firstWallet
		fromWallet = secondWallet
	}

	// Check sufficient balance
	if fromWallet.Balance.LessThan(request.Amount) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Insufficient balance")
	}

	// Create transaction record
	description := request.Description
	transaction := &entity.Transaction{
		Type:              entity.TransactionTypeTransfer,
		Amount:            request.Amount,
		FromWalletID:      &fromWallet.ID,
		ToWalletID:        toWallet.ID,
		PerformedByUserID: *auth.UserID,
		Status:            entity.TransactionStatusCompleted,
		Description:       &description,
	}

	if err := uc.TransactionRepository.Create(tx, transaction); err != nil {
		uc.Log.Errorf("Transaction creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create debit mutation for sender
	fromBalanceBefore := fromWallet.Balance
	fromBalanceAfter := fromWallet.Balance.Sub(request.Amount)

	debitMutation := &entity.WalletMutation{
		WalletID:      fromWallet.ID,
		TransactionID: transaction.ID,
		Type:          entity.MutationTypeDebit,
		Amount:        request.Amount,
		BalanceBefore: fromBalanceBefore,
		BalanceAfter:  fromBalanceAfter,
	}

	if err := uc.WalletMutationRepository.Create(tx, debitMutation); err != nil {
		uc.Log.Errorf("Debit mutation creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create credit mutation for recipient
	toBalanceBefore := toWallet.Balance
	toBalanceAfter := toWallet.Balance.Add(request.Amount)

	creditMutation := &entity.WalletMutation{
		WalletID:      toWallet.ID,
		TransactionID: transaction.ID,
		Type:          entity.MutationTypeCredit,
		Amount:        request.Amount,
		BalanceBefore: toBalanceBefore,
		BalanceAfter:  toBalanceAfter,
	}

	if err := uc.WalletMutationRepository.Create(tx, creditMutation); err != nil {
		uc.Log.Errorf("Credit mutation creation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Update sender wallet balance
	if err := uc.WalletRepository.UpdateBalance(tx, fromWallet.ID, fromBalanceAfter); err != nil {
		uc.Log.Errorf("UpdateBalance error for sender: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Update recipient wallet balance
	if err := uc.WalletRepository.UpdateBalance(tx, toWallet.ID, toBalanceAfter); err != nil {
		uc.Log.Errorf("UpdateBalance error for recipient: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		uc.Log.Errorf("Transaction commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Send real-time notifications
	if uc.Notifier != nil {
		go func() {
			// Notification for sender (confirmation)
			senderNotification := &model.TransactionNotification{
				TransactionID:     transaction.ID,
				TransactionType:   string(transaction.Type),
				Amount:            request.Amount.String(),
				FromUserID:        auth.UserID,
				ToUserID:          request.ToUserID,
				PerformedByUserID: *auth.UserID,
				Description:       transaction.Description,
				CreatedAt:         time.Now().Format(time.RFC3339),
			}
			uc.Notifier.NotifyTransaction(*auth.UserID, senderNotification)

			// Wallet update for sender
			senderWalletNotification := &model.WalletUpdateNotification{
				WalletID:      fromWallet.ID,
				NewBalance:    fromBalanceAfter.String(),
				MutationType:  string(entity.MutationTypeDebit),
				MutationID:    debitMutation.ID,
				TransactionID: transaction.ID,
				Amount:        request.Amount.String(),
			}
			uc.Notifier.NotifyWalletUpdate(*auth.UserID, senderWalletNotification)

			// Notification for recipient
			recipientNotification := &model.TransactionNotification{
				TransactionID:     transaction.ID,
				TransactionType:   string(transaction.Type),
				Amount:            request.Amount.String(),
				FromUserID:        auth.UserID,
				ToUserID:          request.ToUserID,
				PerformedByUserID: *auth.UserID,
				Description:       transaction.Description,
				CreatedAt:         time.Now().Format(time.RFC3339),
			}
			uc.Notifier.NotifyTransaction(request.ToUserID, recipientNotification)

			// Wallet update for recipient
			recipientWalletNotification := &model.WalletUpdateNotification{
				WalletID:      toWallet.ID,
				NewBalance:    toBalanceAfter.String(),
				MutationType:  string(entity.MutationTypeCredit),
				MutationID:    creditMutation.ID,
				TransactionID: transaction.ID,
				Amount:        request.Amount.String(),
			}
			uc.Notifier.NotifyWalletUpdate(request.ToUserID, recipientWalletNotification)
		}()
	}

	return converter.TransactionToTransactionResponse(transaction), nil
}

// GetTransactionsByUserID retrieves transactions for a user.
func (uc *TransactionUseCase) GetTransactionsByUserID(ctx context.Context, userID uint, page, limit int) (*model.TransactionListResponse, error) {
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

	transactions, total, err := uc.TransactionRepository.FindByWalletID(uc.DB.WithContext(ctx), wallet.ID, page, limit)
	if err != nil {
		uc.Log.Errorf("FindByWalletID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.TransactionListResponse{
		Transactions: converter.TransactionsToTransactionResponses(transactions),
		Total:        total,
		Page:         page,
		Limit:        limit,
	}, nil
}
