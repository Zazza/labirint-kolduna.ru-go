package validation

import (
	"gamebook-backend/modules/auth/dto"
	userDto "gamebook-backend/modules/user/dto"
	"github.com/go-playground/validator/v10"
)

type AuthValidation struct {
	validate *validator.Validate
}

func NewAuthValidation() *AuthValidation {
	validate := validator.New()

	validate.RegisterValidation("password", validatePassword)

	return &AuthValidation{
		validate: validate,
	}
}

func (v *AuthValidation) ValidateRegisterRequest(req userDto.UserCreateRequest) error {
	return v.validate.Struct(req)
}

func (v *AuthValidation) ValidateLoginRequest(req userDto.UserLoginRequest) error {
	return v.validate.Struct(req)
}

func (v *AuthValidation) ValidateRefreshTokenRequest(req dto.RefreshTokenRequest) error {
	return v.validate.Struct(req)
}

// Custom validators
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	// Add more password validation rules as needed
	return true
}
