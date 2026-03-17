package validation

import (
	"errors"
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
	err := v.validate.Struct(req)
	if err != nil {
		return err
	}

	return ValidatePasswordComplexity(req.Password)
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

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= 32 && char <= 126:
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return false
	}

	return true
}

func ValidatePasswordComplexity(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= 32 && char <= 126:
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
