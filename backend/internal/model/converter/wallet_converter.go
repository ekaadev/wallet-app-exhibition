package converter

import (
	"backend/internal/entity"
	"backend/internal/model"
)

func WalletToWalletResponse(wallet *entity.Wallet) *model.WalletResponse {
	return &model.WalletResponse{
		ID:      wallet.ID,
		UserID:  wallet.UserID,
		Balance: wallet.Balance,
	}
}
