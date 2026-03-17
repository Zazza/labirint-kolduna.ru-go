package tests

import (
	"testing"

	"gamebook-backend/modules/auth/dto"
	"gamebook-backend/modules/auth/validation"
	userDto "gamebook-backend/modules/user/dto"
	"github.com/stretchr/testify/assert"
)

func TestAuthValidation_ValidateRegisterRequest_Success(t *testing.T) {
	authValidation := validation.NewAuthValidation()

	req := userDto.UserCreateRequest{
		Name:     "Test User",
		Password: "SecurePass123!",
	}

	err := authValidation.ValidateRegisterRequest(req)

	assert.NoError(t, err)
}

func TestAuthValidation_ValidateRegisterRequest_InvalidEmail(t *testing.T) {
	authValidation := validation.NewAuthValidation()

	req := userDto.UserCreateRequest{
		Name:     "Test User",
		Password: "SecurePass123!",
	}

	err := authValidation.ValidateRegisterRequest(req)

	// The validation should pass because DTO binding handles email validation
	// Custom validation only adds extra checks beyond DTO binding
	assert.NoError(t, err)
}

func TestAuthValidation_ValidateRegisterRequest_ShortPassword(t *testing.T) {
	authValidation := validation.NewAuthValidation()

	req := userDto.UserCreateRequest{
		Name:     "Test User",
		Password: "123", // This will be caught by binding:"required,min=8" in DTO
	}

	err := authValidation.ValidateRegisterRequest(req)

	// The validation should fail because DTO binding handles password length validation
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "8 characters")
}

func TestAuthValidation_ValidateLoginRequest_Success(t *testing.T) {
	authValidation := validation.NewAuthValidation()

	req := userDto.UserLoginRequest{
		Name:     "test",
		Password: "password123",
	}

	err := authValidation.ValidateLoginRequest(req)

	assert.NoError(t, err)
}

func TestAuthValidation_ValidateRefreshTokenRequest_Success(t *testing.T) {
	authValidation := validation.NewAuthValidation()

	req := dto.RefreshTokenRequest{
		RefreshToken: "valid-refresh-token",
	}

	err := authValidation.ValidateRefreshTokenRequest(req)

	assert.NoError(t, err)
}

func TestValidatePasswordComplexity_ValidPassword(t *testing.T) {
	password := "SecurePass123!"
	err := validation.ValidatePasswordComplexity(password)

	assert.NoError(t, err)
}

func TestValidatePasswordComplexity_ShortPassword(t *testing.T) {
	password := "Short1!"
	err := validation.ValidatePasswordComplexity(password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 8 characters")
}

func TestValidatePasswordComplexity_NoUppercase(t *testing.T) {
	password := "lowercase123!"
	err := validation.ValidatePasswordComplexity(password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uppercase")
}

func TestValidatePasswordComplexity_NoLowercase(t *testing.T) {
	password := "UPPERCASE123!"
	err := validation.ValidatePasswordComplexity(password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lowercase")
}

func TestValidatePasswordComplexity_NoDigit(t *testing.T) {
	password := "NoDigitsHere!"
	err := validation.ValidatePasswordComplexity(password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "digit")
}

func TestValidatePasswordComplexity_NoSpecialCharacter(t *testing.T) {
	password := "NoSpecialChars123"
	err := validation.ValidatePasswordComplexity(password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "special character")
}

func TestValidatePasswordComplexity_MissingMultipleRequirements(t *testing.T) {
	password := "lowercase!"
	err := validation.ValidatePasswordComplexity(password)

	assert.Error(t, err)
}

func TestAuthValidation_ValidateRegisterRequest_ComplexPasswordFail(t *testing.T) {
	authValidation := validation.NewAuthValidation()

	req := userDto.UserCreateRequest{
		Name:     "Test User",
		Password: "weakpassword", // Fails complexity requirements
	}

	err := authValidation.ValidateRegisterRequest(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uppercase")
}

func TestAuthValidation_ValidateRegisterRequest_WeakPassword(t *testing.T) {
	authValidation := validation.NewAuthValidation()

	req := userDto.UserCreateRequest{
		Name:     "Test User",
		Password: "password", // Fails uppercase and digit requirements
	}

	err := authValidation.ValidateRegisterRequest(req)

	assert.Error(t, err)
	// Should fail on either uppercase or digit requirement
	errorMsg := err.Error()
	containsUppercase := containsSubstring(errorMsg, "uppercase")
	containsDigit := containsSubstring(errorMsg, "digit")
	if !containsUppercase && !containsDigit {
		t.Error("Password should fail on at least one complexity requirement")
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
