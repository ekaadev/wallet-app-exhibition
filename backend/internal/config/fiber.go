package config

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// NewFiber creates and returns a new Fiber app instance configured for the application.
// Fiber app for handling HTTP requests.
func NewFiber(viper *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      viper.GetString("app.name"),
		Prefork:      viper.GetBool("web.prefork"),
		ErrorHandler: NewErrorHandler(),
	})

	return app
}

// NewErrorHandler creates and returns a new Fiber error handler for the application.
// Error handler for handling errors in Fiber app.
func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		// Default error
		code := fiber.StatusInternalServerError

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
