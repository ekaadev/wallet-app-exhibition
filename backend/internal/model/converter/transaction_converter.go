package converter

import (
	"backend/internal/entity"
	"backend/internal/model"
)

func TransactionToTransactionResponse(transaction *entity.Transaction) *model.TransactionResponse {
	return &model.TransactionResponse{
		ID:                transaction.ID,
		Type:              string(transaction.Type),
		Amount:            transaction.Amount,
		FromWalletID:      transaction.FromWalletID,
		ToWalletID:        transaction.ToWalletID,
		PerformedByUserID: transaction.PerformedByUserID,
		Status:            string(transaction.Status),
		Description:       transaction.Description,
		CreatedAt:         transaction.CreatedAt,
	}
}

func TransactionsToTransactionResponses(transactions []entity.Transaction) []model.TransactionResponse {
	responses := make([]model.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *TransactionToTransactionResponse(&transaction)
	}
	return responses
}
