package http

import (
	"backend/internal/delivery/http/middleware"
	"backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WalletController struct {
	Log           *logrus.Logger
	WalletUseCase usecase.WalletUseCaseInterface
}

func NewWalletController(log *logrus.Logger, walletUseCase usecase.WalletUseCaseInterface) *WalletController {
	return &WalletController{
		Log:           log,
		WalletUseCase: walletUseCase,
	}
}
func (wc *WalletController) GetMyWallet(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	response, err := wc.WalletUseCase.GetByUserID(ctx.UserContext(), *auth.UserID)
	if err != nil {
		wc.Log.Warnf("WalletUseCase.GetByUserID error: %v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}
