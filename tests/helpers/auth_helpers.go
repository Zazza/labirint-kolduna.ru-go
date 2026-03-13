package helpers

import (
	authService "gamebook-backend/modules/auth/service"

	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func GenerateTestJWT(t *testing.T, userID uuid.UUID, jwtSvc authService.JWTService) string {
	token := jwtSvc.GenerateAccessToken(userID.String())
	require.NotEmpty(t, token)
	return token
}

func SetupTestAuthService(t *testing.T) authService.JWTService {
	jwtSvc := authService.NewJWTService()
	require.NotNil(t, jwtSvc)
	return jwtSvc
}

func CreateAuthenticatedRequest(t *testing.T, userID uuid.UUID, jwtSvc authService.JWTService) string {
	token := GenerateTestJWT(t, userID, jwtSvc)
	require.NotEmpty(t, token)
	return token
}
