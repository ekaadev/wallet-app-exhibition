package http

import (
	"backend/internal/delivery/http/middleware"
	"backend/internal/model"
	"backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log         *logrus.Logger
	UserUseCase usecase.UserUseCaseInterface
}

func NewUserController(log *logrus.Logger, userUseCase usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		Log:         log,
		UserUseCase: userUseCase,
	}
}

func (uc *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.UserRegistrationRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		uc.Log.Warnf("BodyParser error: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := uc.UserUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		uc.Log.Warnf("UserUseCase.Create error: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": response,
	})
}

func (uc *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.UserLoginRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		uc.Log.Warnf("BodyParser error: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := uc.UserUseCase.Login(ctx.UserContext(), request)
	if err != nil {
		uc.Log.Warnf("UserUseCase.Login error: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

// GetProfile returns the authenticated user's profile with wallet information.
func (uc *UserController) GetProfile(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	response, err := uc.UserUseCase.GetProfile(ctx.UserContext(), *auth.UserID)
	if err != nil {
		uc.Log.Warnf("UserUseCase.GetProfile error: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}
