package routes

import (
	"gamebook-backend/database/entities"
	gameDto "gamebook-backend/modules/game/dto"
	helpers2 "gamebook-backend/tests/helpers"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChoice_Success(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)

	targetSection := helpers2.CreateTestSection(t, testCtx.DB, 2, entities.SectionType("normal"))
	transition := helpers2.CreateTestTransition(t, testCtx.DB, section.ID, targetSection.ID, "Go to section 2")

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	reqBody := map[string]interface{}{
		"transitionID": transition.ID,
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/choice", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.ActionResponse](t, resp)
	assert.True(t, apiResp.Success)
}

func TestChoice_Unauthorized(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	reqBody := map[string]interface{}{
		"transitionID": uuid.New(),
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/choice", "", reqBody)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestChoice_ChangesSection(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 20, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)

	oldSectionID := player.SectionID

	targetSection := helpers2.CreateTestSection(t, testCtx.DB, 21, entities.SectionType("normal"))
	transition := helpers2.CreateTestTransition(t, testCtx.DB, section.ID, targetSection.ID, "Go to section 21")

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	reqBody := map[string]interface{}{
		"transitionID": transition.ID,
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/choice", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedPlayer entities.Player
	err := testCtx.DB.Preload("Section").First(&updatedPlayer, player.ID).Error
	require.NoError(t, err)

	assert.NotEqual(t, oldSectionID, updatedPlayer.SectionID, "Player section ID should change after choice")
	assert.Equal(t, targetSection.ID, updatedPlayer.SectionID, "Player should be moved to target section")
}
