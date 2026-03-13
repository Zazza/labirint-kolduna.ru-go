package routes

import (
	"gamebook-backend/database/entities"
	gameDto "gamebook-backend/modules/game/dto"
	helpers2 "gamebook-backend/tests/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRollTheDice_Success(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	resp := helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/role-the-dice", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.RollTheDiceDto](t, resp)
	assert.True(t, apiResp.Success)
}

func TestRollTheDice_Unauthorized(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	resp := helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/role-the-dice", "", nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestRollTheDice_ValidRange(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 10, entities.SectionType("normal"))
	helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	resp := helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/role-the-dice", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := helpers2.ParseData[gameDto.RollTheDiceDto](t, resp)
	assert.GreaterOrEqual(t, result.DiceFirst, uint(1))
	assert.LessOrEqual(t, result.DiceFirst, uint(6))
	assert.GreaterOrEqual(t, result.DiceSecond, uint(1))
	assert.LessOrEqual(t, result.DiceSecond, uint(6))
}
