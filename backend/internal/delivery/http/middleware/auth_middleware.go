package middleware

import (
	"backend/internal/model"
	"backend/internal/usecase"
	"backend/internal/util"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// NewAuth creates a new authentication middleware for Fiber.
func NewAuth(userUseCase *usecase.UserUseCase, tokenUtil *util.TokenUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the Authorization header
		request := &model.VerifyUserRequest{
			Token: ctx.Get("Authorization", "NOT_FOUND"),
		}

		// Split the "Bearer" prefix from the token
		parts := strings.Split(request.Token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			userUseCase.Log.Warnf("Invalid authorization header format")
			return fiber.ErrUnauthorized
		}

		// Parse and verify the token
		auth, err := tokenUtil.ParseToken(ctx.UserContext(), parts[1])
		if err != nil {
			userUseCase.Log.Warnf("Invalid authorization header format")
			return fiber.ErrUnauthorized
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
