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
		var tokenString string

		// 1. Coba ambil dari Header Authorization
		authHeader := ctx.Get("Authorization")

		if authHeader != "" {
			// Format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// 2. Jika di Header tidak ada, coba ambil dari Cookies
		if tokenString == "" {
			tokenString = ctx.Cookies("jwt")
		}

		// 3. Jika di Cookies tidak ada, coba ambil dari Query Param "token" (untuk WebSocket)
		if tokenString == "" {
			tokenString = ctx.Query("token")
		}

		// 3. Jika tetap kosong, baru return Unauthorized
		if tokenString == "" {
			userUseCase.Log.Warnf("Missing or invalid authorization format")
			return fiber.ErrUnauthorized
		}

		// 4. Parse dan verify token yang didapat
		auth, err := tokenUtil.ParseToken(ctx.UserContext(), tokenString)
		if err != nil {
			userUseCase.Log.Warnf("Invalid token verification: %v", err)
			return fiber.ErrUnauthorized
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
