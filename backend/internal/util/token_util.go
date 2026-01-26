package util

import (
	"backend/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type TokenUtil struct {
	SecretKey string
	Redis     *redis.Client
}

// NewTokenUtil creates a new instance of TokenUtil.
// TokenUtil requires a secret key for signing tokens and a Redis client for token storage and management.
func NewTokenUtil(secretKey string, redisClient *redis.Client) *TokenUtil {
	return &TokenUtil{
		SecretKey: secretKey,
		Redis:     redisClient,
	}
}

// CreateToken generates a new JWT token for the given Auth model.
func (tu *TokenUtil) CreateToken(ctx context.Context, auth *model.Auth) (string, error) {
	// Set token now and expiry with in 30 days
	now := time.Now()
	expiryDuration := time.Hour * 24 * 30

	// Assign registered claims
	auth.ExpiresAt = jwt.NewNumericDate(now.Add(expiryDuration))
	auth.IssuedAt = jwt.NewNumericDate(now)

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth)

	// Sign token with secret key
	signedString, err := token.SignedString([]byte(tu.SecretKey))
	if err != nil {
		return "", err
	}

	// Save token in Redis with expiry
	_, err = tu.Redis.Set(ctx, signedString, "valid", expiryDuration).Result()
	if err != nil {
		return "", err
	}

	return signedString, nil
}

// ParseToken validates and parses the given JWT token string.
func (tu *TokenUtil) ParseToken(ctx context.Context, tokenString string) (*model.Auth, error) {
	// Parse token with claims, receiving jwt structure
	token, err := jwt.ParseWithClaims(tokenString, &model.Auth{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tu.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	claims, ok := token.Claims.(*model.Auth)
	if !ok || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	// Check token expiry
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fiber.ErrUnauthorized
	}

	// Check token existence in Redis
	result, err := tu.Redis.Exists(ctx, tokenString).Result()
	if err != nil {
		return nil, err
	}

	if result == 0 {
		return nil, fiber.ErrUnauthorized
	}

	return claims, nil
}
