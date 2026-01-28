package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// NewViper creates and returns a new Viper instance configured for the application.
// Viper for configuration management.
func NewViper() *viper.Viper {
	config := viper.New()
	config.AutomaticEnv()

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config.SetDefault("app.name", "wallet-exhition-app")
	config.SetDefault("web.port", 3000)
	config.SetDefault("cookie.secure", false)
	config.SetDefault("log.level", 7)

	config.SetDefault("database.pool.idle", 10)
	config.SetDefault("database.pool.max", 100)
	config.SetDefault("database.pool.lifetime", 300)

	// Read config.json
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./backend")
	config.AddConfigPath("../../backend")
	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		//panic(fmt.Errorf("Fatal error config file: %s \n", err))
		fmt.Println("Skipping! Failed to read config:", err)
	}

	// Read .env file and merge with existing config
	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath("./backend")
	config.AddConfigPath("../../backend")
	config.AddConfigPath(".")
	if err := config.MergeInConfig(); err != nil {
		//panic(fmt.Errorf("Fatal error env file: %s \n", err))
		fmt.Println("Skipping! Failed to read env file:", err)
	}

	return config
}
