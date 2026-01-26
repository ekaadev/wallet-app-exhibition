package controller_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpDelivery "backend/internal/delivery/http"
	"backend/internal/model"
	"backend/tests/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupWalletMutationTestApp creates a Fiber app with WalletMutationController for testing.
func setupWalletMutationTestApp(mockUseCase *mocks.MockWalletMutationUseCase) *fiber.App {
	app := fiber.New()
	log := logrus.New()
	log.SetOutput(io.Discard)

	controller := httpDelivery.NewWalletMutationController(log, mockUseCase)

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

	app.Get("/wallet-mutations", controller.GetMyMutations)

	return app
}

// TestGetMyMutations_Success tests successful wallet mutations retrieval.
func TestGetMyMutations_Success(t *testing.T) {
	mockUseCase := new(mocks.MockWalletMutationUseCase)
	app := setupWalletMutationTestApp(mockUseCase)

	expectedResponse := &model.WalletMutationListResponse{
		Mutations: []model.WalletMutationResponse{
			{
				ID:            1,
				WalletID:      1,
				TransactionID: 1,
				Type:          "credit",
				Amount:        decimal.NewFromInt(100000),
				BalanceBefore: decimal.NewFromInt(0),
				BalanceAfter:  decimal.NewFromInt(100000),
				CreatedAt:     time.Now(),
			},
		},
		Total: 1,
		Page:  1,
		Limit: 10,
	}

	mockUseCase.On("GetMutationsByUserID", mock.Anything, uint(1), 1, 10).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/wallet-mutations", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
	assert.Equal(t, float64(1), data["page"])
	assert.Equal(t, float64(10), data["limit"])

	mutations := data["mutations"].([]interface{})
	assert.Len(t, mutations, 1)

	mockUseCase.AssertExpectations(t)
}

// TestGetMyMutations_WithPagination tests wallet mutations with pagination parameters.
func TestGetMyMutations_WithPagination(t *testing.T) {
	mockUseCase := new(mocks.MockWalletMutationUseCase)
	app := setupWalletMutationTestApp(mockUseCase)

	expectedResponse := &model.WalletMutationListResponse{
		Mutations: []model.WalletMutationResponse{},
		Total:     0,
		Page:      2,
		Limit:     5,
	}

	mockUseCase.On("GetMutationsByUserID", mock.Anything, uint(1), 2, 5).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/wallet-mutations?page=2&limit=5", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["page"])
	assert.Equal(t, float64(5), data["limit"])

	mockUseCase.AssertExpectations(t)
}

// TestGetMyMutations_WalletNotFound tests wallet mutations when wallet is not found.
func TestGetMyMutations_WalletNotFound(t *testing.T) {
	mockUseCase := new(mocks.MockWalletMutationUseCase)
	app := setupWalletMutationTestApp(mockUseCase)

	mockUseCase.On("GetMutationsByUserID", mock.Anything, uint(1), 1, 10).
		Return(nil, fiber.NewError(fiber.StatusNotFound, "Wallet not found"))

	req := httptest.NewRequest(http.MethodGet, "/wallet-mutations", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockUseCase.AssertExpectations(t)
}
