package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewViper creates and returns a new Viper instance configured for the application.
// Viper for configuration management.
func NewViper() *viper.Viper {
	config := viper.New()

	// Read config.json
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./backend")
	config.AddConfigPath("../../backend")
	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Read .env file and merge with existing config
	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath("./backend")
	config.AddConfigPath("../../backend")
	config.AddConfigPath(".")
	config.AutomaticEnv()
	if err := config.MergeInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error env file: %s \n", err))
	}

	return config
}
