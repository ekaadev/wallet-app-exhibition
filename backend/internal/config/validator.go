package config

import "github.com/go-playground/validator/v10"

// NewValidator creates and returns a new Validator instance for input validation.
// Validator for input validation.
func NewValidator() *validator.Validate {
	return validator.New()
}
