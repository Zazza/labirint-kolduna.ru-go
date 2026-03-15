package routes

import (
	"gamebook-backend/database/entities"
	gameDto "gamebook-backend/modules/game/dto"
	helpers2 "gamebook-backend/tests/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSleep_Success(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)

	for i := 202; i <= 212; i++ {
		helpers2.CreateTestSection(t, testCtx.DB, uint(i), entities.SectionType("sleep"))
	}

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/sleep", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)
}

func TestSleep_Unauthorized(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/sleep", "", nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestSleep_DeadPlayer(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)

	player.Health = 0
	testCtx.DB.Save(player)

	for i := 202; i <= 212; i++ {
		helpers2.CreateTestSection(t, testCtx.DB, uint(i), entities.SectionType("sleep"))
	}

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/sleep", token, nil)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	result := helpers2.ParseData[gameDto.SleepDTO](t, resp)

	assert.Equal(t, result.Result, gameDto.ActionResult(false))
}
