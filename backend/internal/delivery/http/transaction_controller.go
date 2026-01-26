package http

import (
	"backend/internal/delivery/http/middleware"
	"backend/internal/model"
	"backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	Log                *logrus.Logger
	TransactionUseCase usecase.TransactionUseCaseInterface
}

func NewTransactionController(log *logrus.Logger, transactionUseCase usecase.TransactionUseCaseInterface) *TransactionController {
	return &TransactionController{
		Log:                log,
		TransactionUseCase: transactionUseCase,
	}
}
func (tc *TransactionController) TopUp(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.TopUpRequest)
	if err := ctx.BodyParser(request); err != nil {
		tc.Log.Warnf("BodyParser error: %v", err)
		return fiber.ErrBadRequest
	}
	response, err := tc.TransactionUseCase.TopUp(ctx.UserContext(), auth, request)
	if err != nil {
		tc.Log.Warnf("TransactionUseCase.TopUp error: %v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": response,
	})
}
func (tc *TransactionController) Transfer(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.TransferRequest)
	if err := ctx.BodyParser(request); err != nil {
		tc.Log.Warnf("BodyParser error: %v", err)
		return fiber.ErrBadRequest
	}
	response, err := tc.TransactionUseCase.Transfer(ctx.UserContext(), auth, request)
	if err != nil {
		tc.Log.Warnf("TransactionUseCase.Transfer error: %v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": response,
	})
}
func (tc *TransactionController) GetMyTransactions(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	response, err := tc.TransactionUseCase.GetTransactionsByUserID(ctx.UserContext(), *auth.UserID, page, limit)
	if err != nil {
		tc.Log.Warnf("TransactionUseCase.GetTransactionsByUserID error: %v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}
