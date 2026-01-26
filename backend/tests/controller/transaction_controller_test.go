package controller_test

import (
	"bytes"
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

// setupTransactionTestApp creates a Fiber app with TransactionController for testing.
func setupTransactionTestApp(mockUseCase *mocks.MockTransactionUseCase, role string) *fiber.App {
	app := fiber.New()
	log := logrus.New()
	log.SetOutput(io.Discard)

	controller := httpDelivery.NewTransactionController(log, mockUseCase)

	// Middleware to set auth context for testing
	app.Use(func(c *fiber.Ctx) error {
		userID := uint(1)
		auth := &model.Auth{
			UserID:   &userID,
			Username: "testuser",
			Role:     role,
		}
		c.Locals("auth", auth)
		return c.Next()
	})

	app.Post("/transactions/topup", controller.TopUp)
	app.Post("/transactions/transfer", controller.Transfer)
	app.Get("/transactions", controller.GetMyTransactions)

	return app
}

// TestTopUp_Success tests successful top-up operation by super admin.
func TestTopUp_Success(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	app := setupTransactionTestApp(mockUseCase, "super_admin")

	toWalletID := uint(2)
	expectedResponse := &model.TransactionResponse{
		ID:                1,
		Type:              "top_up",
		Amount:            decimal.NewFromInt(100000),
		ToWalletID:        toWalletID,
		PerformedByUserID: 1,
		Status:            "completed",
		CreatedAt:         time.Now(),
	}

	mockUseCase.On("TopUp", mock.Anything, mock.Anything, mock.MatchedBy(func(req *model.TopUpRequest) bool {
		return req.ToUserID == 2 && req.Amount.Equal(decimal.NewFromInt(100000))
	})).Return(expectedResponse, nil)

	reqBody := map[string]interface{}{
		"to_user_id": 2,
		"amount":     100000,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions/topup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, "top_up", data["type"])
	assert.Equal(t, "completed", data["status"])

	mockUseCase.AssertExpectations(t)
}

// TestTopUp_InvalidRequest tests top-up with invalid request body.
func TestTopUp_InvalidRequest(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	app := setupTransactionTestApp(mockUseCase, "super_admin")

	req := httptest.NewRequest(http.MethodPost, "/transactions/topup", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestTopUp_Forbidden tests top-up by non-super admin user.
func TestTopUp_Forbidden(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	app := setupTransactionTestApp(mockUseCase, "user")

	mockUseCase.On("TopUp", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, fiber.NewError(fiber.StatusForbidden, "Only super admin can perform top-up"))

	reqBody := map[string]interface{}{
		"to_user_id": 2,
		"amount":     100000,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions/topup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	mockUseCase.AssertExpectations(t)
}

// TestTransfer_Success tests successful transfer operation.
func TestTransfer_Success(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	app := setupTransactionTestApp(mockUseCase, "user")

	fromWalletID := uint(1)
	toWalletID := uint(2)
	expectedResponse := &model.TransactionResponse{
		ID:                1,
		Type:              "transfer",
		Amount:            decimal.NewFromInt(50000),
		FromWalletID:      &fromWalletID,
		ToWalletID:        toWalletID,
		PerformedByUserID: 1,
		Status:            "completed",
		CreatedAt:         time.Now(),
	}

	mockUseCase.On("Transfer", mock.Anything, mock.Anything, mock.MatchedBy(func(req *model.TransferRequest) bool {
		return req.ToUserID == 2 && req.Amount.Equal(decimal.NewFromInt(50000))
	})).Return(expectedResponse, nil)

	reqBody := map[string]interface{}{
		"to_user_id": 2,
		"amount":     50000,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions/transfer", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, "transfer", data["type"])
	assert.Equal(t, "completed", data["status"])

	mockUseCase.AssertExpectations(t)
}

// TestTransfer_InvalidRequest tests transfer with invalid request body.
func TestTransfer_InvalidRequest(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	app := setupTransactionTestApp(mockUseCase, "user")

	req := httptest.NewRequest(http.MethodPost, "/transactions/transfer", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestTransfer_InsufficientBalance tests transfer when balance is insufficient.
func TestTransfer_InsufficientBalance(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	app := setupTransactionTestApp(mockUseCase, "user")

	mockUseCase.On("Transfer", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, fiber.NewError(fiber.StatusBadRequest, "Insufficient balance"))

	reqBody := map[string]interface{}{
		"to_user_id": 2,
		"amount":     1000000,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions/transfer", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	mockUseCase.AssertExpectations(t)
}

// TestGetMyTransactions_Success tests successful transaction list retrieval.
func TestGetMyTransactions_Success(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	app := setupTransactionTestApp(mockUseCase, "user")

	expectedResponse := &model.TransactionListResponse{
		Transactions: []model.TransactionResponse{
			{
				ID:                1,
				Type:              "transfer",
				Amount:            decimal.NewFromInt(50000),
				PerformedByUserID: 1,
				Status:            "completed",
				CreatedAt:         time.Now(),
			},
		},
		Total: 1,
		Page:  1,
		Limit: 10,
	}

	mockUseCase.On("GetTransactionsByUserID", mock.Anything, uint(1), 1, 10).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
	assert.Equal(t, float64(1), data["page"])
	assert.Equal(t, float64(10), data["limit"])

	mockUseCase.AssertExpectations(t)
}

// TestGetMyTransactions_WithPagination tests transaction list with pagination parameters.
func TestGetMyTransactions_WithPagination(t *testing.T) {
	mockUseCase := new(mocks.MockTransactionUseCase)
	app := setupTransactionTestApp(mockUseCase, "user")

	expectedResponse := &model.TransactionListResponse{
		Transactions: []model.TransactionResponse{},
		Total:        0,
		Page:         2,
		Limit:        5,
	}

	mockUseCase.On("GetTransactionsByUserID", mock.Anything, uint(1), 2, 5).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/transactions?page=2&limit=5", nil)

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
