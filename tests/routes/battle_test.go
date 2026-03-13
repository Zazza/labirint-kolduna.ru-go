package routes

import (
	helpers2 "gamebook-backend/tests/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBattle_Unauthorized(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	reqBody := map[string]interface{}{
		"transitionID": "uuid",
		"weapon":       "Sword",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/battle", "", reqBody)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
