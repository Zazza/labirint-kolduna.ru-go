package validation

import (
	"github.com/go-playground/validator/v10"
)

type GameValidation struct {
	validate *validator.Validate
}

func NewAuthValidation() *GameValidation {
	validate := validator.New()

	return &GameValidation{
		validate: validate,
	}
}
