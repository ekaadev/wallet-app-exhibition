package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// NewRedisClient creates and returns a new Redis client instance configured for the application.
// Redis client for interacting with Redis database.
func NewRedisClient(viper *viper.Viper) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", viper.GetString("REDIS_HOST"), viper.GetInt("REDIS_PORT")),
		DB:   viper.GetInt("REDIS_DB"),
	})
}
