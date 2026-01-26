package converter

import (
	"backend/internal/entity"
	"backend/internal/model"
)

func WalletMutationToWalletMutationResponse(mutation *entity.WalletMutation) *model.WalletMutationResponse {
	return &model.WalletMutationResponse{
		ID:            mutation.ID,
		WalletID:      mutation.WalletID,
		TransactionID: mutation.TransactionID,
		Type:          string(mutation.Type),
		Amount:        mutation.Amount,
		BalanceBefore: mutation.BalanceBefore,
		BalanceAfter:  mutation.BalanceAfter,
		CreatedAt:     mutation.CreatedAt,
	}
}

func WalletMutationsToWalletMutationResponses(mutations []entity.WalletMutation) []model.WalletMutationResponse {
	responses := make([]model.WalletMutationResponse, len(mutations))
	for i, mutation := range mutations {
		responses[i] = *WalletMutationToWalletMutationResponse(&mutation)
	}
	return responses
}
