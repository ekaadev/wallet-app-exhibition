package controller_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httpDelivery "backend/internal/delivery/http"
	"backend/internal/model"
	"backend/tests/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserTestApp(mockUseCase *mocks.MockUserUseCase) *fiber.App {
	app := fiber.New()
	log := logrus.New()
	log.SetOutput(io.Discard)
	config := viper.New()

	controller := httpDelivery.NewUserController(log, config, mockUseCase)

	app.Post("/users/register", controller.Register)
	app.Post("/users/login", controller.Login)

	return app
}

// TestRegister_Success tests successful user registration.
func TestRegister_Success(t *testing.T) {
	mockUseCase := new(mocks.MockUserUseCase)
	app := setupUserTestApp(mockUseCase)

	expectedResponse := &model.UserResponse{
		ID:       1,
		Username: "testuser",
		Token:    "jwt-token-here",
	}

	mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(req *model.UserRegistrationRequest) bool {
		return req.Username == "testuser" && req.Password == "password123"
	})).Return(expectedResponse, nil)

	reqBody := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, "testuser", data["username"])
	assert.Equal(t, "jwt-token-here", data["token"])

	mockUseCase.AssertExpectations(t)
}

// TestRegister_InvalidRequest tests registration with invalid request body.
func TestRegister_InvalidRequest(t *testing.T) {
	mockUseCase := new(mocks.MockUserUseCase)
	app := setupUserTestApp(mockUseCase)

	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestRegister_UsernameExists tests registration when username already exists.
func TestRegister_UsernameExists(t *testing.T) {
	mockUseCase := new(mocks.MockUserUseCase)
	app := setupUserTestApp(mockUseCase)

	mockUseCase.On("Create", mock.Anything, mock.Anything).
		Return(nil, fiber.NewError(fiber.StatusConflict, "Username already exists"))

	reqBody := map[string]string{
		"username": "existinguser",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	mockUseCase.AssertExpectations(t)
}

// TestLogin_Success tests successful user login.
func TestLogin_Success(t *testing.T) {
	mockUseCase := new(mocks.MockUserUseCase)
	app := setupUserTestApp(mockUseCase)

	expectedResponse := &model.UserResponse{
		ID:       1,
		Username: "testuser",
		Token:    "jwt-token-here",
	}

	mockUseCase.On("Login", mock.Anything, mock.MatchedBy(func(req *model.UserLoginRequest) bool {
		return req.Username == "testuser" && req.Password == "password123"
	})).Return(expectedResponse, nil)

	reqBody := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, "testuser", data["username"])
	assert.Equal(t, "jwt-token-here", data["token"])

	mockUseCase.AssertExpectations(t)
}

// TestLogin_InvalidRequest tests login with invalid request body.
func TestLogin_InvalidRequest(t *testing.T) {
	mockUseCase := new(mocks.MockUserUseCase)
	app := setupUserTestApp(mockUseCase)

	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestLogin_InvalidCredentials tests login with wrong username/password.
func TestLogin_InvalidCredentials(t *testing.T) {
	mockUseCase := new(mocks.MockUserUseCase)
	app := setupUserTestApp(mockUseCase)

	mockUseCase.On("Login", mock.Anything, mock.Anything).
		Return(nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password"))

	reqBody := map[string]string{
		"username": "wronguser",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockUseCase.AssertExpectations(t)
}
