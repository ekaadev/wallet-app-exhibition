package http

import (
	"backend/internal/delivery/http/middleware"
	"backend/internal/model"
	"backend/internal/usecase"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UserController struct {
	Log         *logrus.Logger
	Config      *viper.Viper
	UserUseCase usecase.UserUseCaseInterface
}

func NewUserController(log *logrus.Logger, config *viper.Viper, userUseCase usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		Log:         log,
		Config:      config,
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

	// Set cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    response.Token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour), // 1 day
		HTTPOnly: true,
		Secure:   uc.Config.GetBool("cookie.secure"),
		SameSite: "Lax",
	})

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

	// Set cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    response.Token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour), // 1 day
		HTTPOnly: true,
		Secure:   uc.Config.GetBool("cookie.secure"),
		SameSite: "Lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

func (uc *UserController) Logout(ctx *fiber.Ctx) error {
	// Clear cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   uc.Config.GetBool("cookie.secure"),
		SameSite: "Lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
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
