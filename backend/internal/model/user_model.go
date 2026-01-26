package model

import "github.com/shopspring/decimal"

// UserRegistrationRequest represents the payload for user registration.
type UserRegistrationRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

// UserLoginRequest represents the payload for user login.
type UserLoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

// VerifyUserRequest represents the payload for user verification in auth middleware.
type VerifyUserRequest struct {
	Token string `validate:"required,max=255"`
}

// UserResponse represents the response payload for user-related operations. (e.g., registration, login)
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// UserProfileResponse represents the response payload for user profile with wallet information.
type UserProfileResponse struct {
	ID       uint                   `json:"id"`
	Username string                 `json:"username"`
	Role     string                 `json:"role"`
	Wallet   *UserProfileWalletInfo `json:"wallet,omitempty"`
}

// UserProfileWalletInfo represents wallet information in user profile response.
type UserProfileWalletInfo struct {
	ID      uint            `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}
