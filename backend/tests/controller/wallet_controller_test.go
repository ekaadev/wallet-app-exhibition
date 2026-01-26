package controller_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httpDelivery "backend/internal/delivery/http"
	"backend/internal/model"
	"backend/tests/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupWalletTestApp creates a Fiber app with WalletController for testing.
func setupWalletTestApp(mockUseCase *mocks.MockWalletUseCase) *fiber.App {
	app := fiber.New()
	log := logrus.New()
	log.SetOutput(io.Discard)

	controller := httpDelivery.NewWalletController(log, mockUseCase)

	// Middleware to set auth context for testing
	app.Use(func(c *fiber.Ctx) error {
		userID := uint(1)
		auth := &model.Auth{
			UserID:   &userID,
			Username: "testuser",
			Role:     "user",
		}
		c.Locals("auth", auth)
		return c.Next()
	})

	app.Get("/wallets/me", controller.GetMyWallet)

	return app
}

// TestGetMyWallet_Success tests successful wallet retrieval.
func TestGetMyWallet_Success(t *testing.T) {
	mockUseCase := new(mocks.MockWalletUseCase)
	app := setupWalletTestApp(mockUseCase)

	expectedResponse := &model.WalletResponse{
		ID:      1,
		UserID:  1,
		Balance: decimal.NewFromInt(100000),
	}

	mockUseCase.On("GetByUserID", mock.Anything, uint(1)).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/wallets/me", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, float64(1), data["user_id"])

	mockUseCase.AssertExpectations(t)
}

// TestGetMyWallet_NotFound tests wallet retrieval when wallet is not found.
func TestGetMyWallet_NotFound(t *testing.T) {
	mockUseCase := new(mocks.MockWalletUseCase)
	app := setupWalletTestApp(mockUseCase)

	mockUseCase.On("GetByUserID", mock.Anything, uint(1)).
		Return(nil, fiber.NewError(fiber.StatusNotFound, "Wallet not found"))

	req := httptest.NewRequest(http.MethodGet, "/wallets/me", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockUseCase.AssertExpectations(t)
}
