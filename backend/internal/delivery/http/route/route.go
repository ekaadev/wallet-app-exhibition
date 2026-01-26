package route

import (
	"backend/internal/delivery/http"
	"backend/internal/delivery/websocket"

	"github.com/gofiber/fiber/v2"
)

type ConfigRoute struct {
	App                      *fiber.App
	UserController           *http.UserController
	WalletController         *http.WalletController
	TransactionController    *http.TransactionController
	WalletMutationController *http.WalletMutationController
	WebSocketHandler         *websocket.Handler
	AuthMiddleware           fiber.Handler
}

// Setup sets up the main routes for the application.
func (cr *ConfigRoute) Setup() {
	// Health check route
	cr.App.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})
	cr.SetupGuestRoutes()
	cr.SetupAuthRoutes()
	cr.SetupWebSocketRoutes()
}

// SetupGuestRoutes sets up routes accessible to guest users.
func (cr *ConfigRoute) SetupGuestRoutes() {
	// User routes (public)
	cr.App.Post("/users/register", cr.UserController.Register)
	cr.App.Post("/users/login", cr.UserController.Login)
}

func (cr *ConfigRoute) SetupAuthRoutes() {
	// Protected routes
	auth := cr.App.Group("", cr.AuthMiddleware)

	// Wallet routes
	auth.Get("/wallets/me", cr.WalletController.GetMyWallet)

	// Transaction routes
	auth.Post("/transactions/topup", cr.TransactionController.TopUp)
	auth.Post("/transactions/transfer", cr.TransactionController.Transfer)
	auth.Get("/transactions", cr.TransactionController.GetMyTransactions)

	// Wallet Mutation routes
	auth.Get("/wallet-mutations", cr.WalletMutationController.GetMyMutations)
}

// SetupWebSocketRoutes sets up WebSocket routes for real-time features.
func (cr *ConfigRoute) SetupWebSocketRoutes() {
	if cr.WebSocketHandler == nil {
		return
	}

	// WebSocket route with token-based authentication via query parameter
	cr.App.Get("/ws", cr.WebSocketHandler.UpgradeMiddleware(), cr.WebSocketHandler.HandleConnection())
}
