package main

import (
	"backend/internal/config"
	"fmt"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogger(viper)
	validator := config.NewValidator()
	redis := config.NewRedisClient(viper)
	db := config.NewDatabase(viper, log)
	app := config.NewFiber(viper)

	config.Bootstrap(&config.BootstrapConfig{
		App:       app,
		DB:        db,
		Redis:     redis,
		Log:       log,
		Validator: validator,
		Config:    viper,
	})

	webPort := viper.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
