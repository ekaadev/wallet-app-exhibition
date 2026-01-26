package http

import (
	"backend/internal/delivery/http/middleware"
	"backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WalletMutationController struct {
	Log                   *logrus.Logger
	WalletMutationUseCase usecase.WalletMutationUseCaseInterface
}

func NewWalletMutationController(log *logrus.Logger, walletMutationUseCase usecase.WalletMutationUseCaseInterface) *WalletMutationController {
	return &WalletMutationController{
		Log:                   log,
		WalletMutationUseCase: walletMutationUseCase,
	}
}
func (wmc *WalletMutationController) GetMyMutations(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	response, err := wmc.WalletMutationUseCase.GetMutationsByUserID(ctx.UserContext(), *auth.UserID, page, limit)
	if err != nil {
		wmc.Log.Warnf("WalletMutationUseCase.GetMutationsByUserID error: %v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}
