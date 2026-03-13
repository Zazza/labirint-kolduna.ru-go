package routes

import (
	"gamebook-backend/database/entities"
	gameDto "gamebook-backend/modules/game/dto"
	helpers2 "gamebook-backend/tests/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSection_Success(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)

	targetSection := helpers2.CreateTestSection(t, testCtx.DB, 2, entities.SectionType("normal"))
	helpers2.CreateTestTransition(t, testCtx.DB, section.ID, targetSection.ID, "Go to section 2")

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	resp := helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/get-section", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)
}

func TestGetSection_Unauthorized(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	resp := helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/get-section", "", nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetSection_WithTransitions(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 10, entities.SectionType("normal"))
	helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)

	targetSection := helpers2.CreateTestSection(t, testCtx.DB, 11, entities.SectionType("normal"))
	helpers2.CreateTestTransition(t, testCtx.DB, section.ID, targetSection.ID, "Option 1")
	helpers2.CreateTestTransition(t, testCtx.DB, section.ID, targetSection.ID, "Option 2")

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	resp := helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/get-section", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := helpers2.ParseData[gameDto.CurrentResponse](t, resp)
	assert.NotEmpty(t, result.Transitions)
	assert.Greater(t, len(result.Transitions), 0)
}

func TestGetSection_BattleSection(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)

	battleStart := "dices"
	section, enemies := helpers2.SetupBattleSection(t, testCtx.DB, 30, &battleStart)

	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	helpers2.CreatePlayerSection(t, testCtx.DB, entities.PlayerSection{
		SectionID: section.ID,
		PlayerID:  player.ID,
	})
	for _, enemy := range enemies {
		helpers2.CreatePlayerSectionEnemy(t, testCtx.DB, player.ID, section.ID, enemy.ID)
	}

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	resp := helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/get-section", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := helpers2.ParseData[gameDto.CurrentResponse](t, resp)
	assert.Equal(t, true, result.RollTheDices)

	resp = helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/role-the-dice", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp = helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/get-section", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result = helpers2.ParseData[gameDto.CurrentResponse](t, resp)
	assert.Equal(t, gameDto.SectionTypeBattle, result.Type)
	assert.Greater(t, len(result.Transitions), 0)
}
