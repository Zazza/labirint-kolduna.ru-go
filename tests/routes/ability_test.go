package routes

import (
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bonus/helpers"
	"gamebook-backend/modules/game/dice"
	gameDto "gamebook-backend/modules/game/dto"
	helpers2 "gamebook-backend/tests/helpers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeds_Success(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	player.Health = 1
	player.Meds = entities.Meds{Name: "Лекарства", Item: "chain mail", Count: 15}
	testCtx.DB.Save(player)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/meds", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.Equal(t, player.Meds.Count, 14)
}

func TestMeds_PlayerDead(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	player.Health = 0
	player.Meds = entities.Meds{Name: "Лекарства", Item: "chain mail", Count: 15}
	testCtx.DB.Save(player)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/meds", token, nil)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.False(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.Equal(t, player.Meds.Count, 15)
}

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

func TestBonus_LuckyStone(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	helpers2.CreatePlayerSection(t, testCtx.DB, entities.PlayerSection{
		SectionID: section.ID,
		PlayerID:  player.ID,
	})

	player.Bonus = getBonuses("Счастливый камушек", "lucky_stone")
	testCtx.DB.Save(player)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	reqBody := map[string]interface{}{
		"bonus":  "lucky_stone",
		"option": "",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.Equal(t, true, helpers.HasBuff(player.Buff, entities.DebuffAliasLuckyStoneReason))

	service := dice.NewPlayerRollTheDices(testCtx.DB, player)
	successResult := false
	for i := 0; i <= 100; i++ {
		result, _ := service.RollTheDice(t.Context(), *player)
		if 6 < *result {
			successResult = true
			break
		}
	}
	assert.True(t, successResult)
}

func TestBonus_AntiPoisonSpell(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	helpers2.CreatePlayerSection(t, testCtx.DB, entities.PlayerSection{
		SectionID: section.ID,
		PlayerID:  player.ID,
	})

	player.Bonus = getBonuses("Заклинание смерти", "anti_poison_spell")
	health := uint(30)
	player.Debuff = []entities.Debuff{
		{Alias: "poison", Health: &health},
	}
	testCtx.DB.Save(player)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	reqBody := map[string]interface{}{
		"bonus":  "anti_poison_spell",
		"option": "",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.Equal(t, []entities.Debuff(nil), player.Debuff)
}

func TestBonus_DeathSpell(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	battleStart := "player"
	section, enemies := helpers2.SetupBattleSection(t, testCtx.DB, 1, 30, &battleStart)
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

	player.Bonus = getBonuses("Заклинание смерти", "death_spell")
	testCtx.DB.Save(player)

	reqBody := map[string]interface{}{
		"bonus":  "death_spell",
		"option": "",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	resp = helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/get-section", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := helpers2.ParseData[gameDto.CurrentResponse](t, resp)
	assert.Contains(t, result.Text, "Заклинание смерти прикончило врага")
}

func TestBonus_InstantHypnosisSpell(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	battleStart := "player"
	section, enemies := helpers2.SetupBattleSection(t, testCtx.DB, 30, 2, &battleStart)
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

	player.Bonus = getBonuses("Заклинание Мгновенного Гипноза", "instant_hypnosis_spell")
	testCtx.DB.Save(player)

	reqBody := map[string]interface{}{
		"bonus":  "instant_hypnosis_spell",
		"option": "",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp = helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/get-section", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := helpers2.ParseData[gameDto.CurrentResponse](t, resp)
	assert.Contains(t, result.Text, "враг впадает в транс и больше вас не беспокоит")
}

func TestBonus_InstantMovement(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	helpers2.CreatePlayerSection(t, testCtx.DB, entities.PlayerSection{
		SectionID: section.ID,
		PlayerID:  player.ID,
	})

	section3 := helpers2.CreateTestSection(t, testCtx.DB, 3, entities.SectionType("normal"))

	player.Bonus = getBonuses("Заклинание Мгновенного Перемещения", "instant_movement")
	testCtx.DB.Save(player)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	reqBody := map[string]interface{}{
		"bonus":  "instant_movement",
		"option": "",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.Equal(t, section3.ID, player.SectionID)
}

func TestBonus_InstantRecoverySpell(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	helpers2.CreatePlayerSection(t, testCtx.DB, entities.PlayerSection{
		SectionID: section.ID,
		PlayerID:  player.ID,
	})

	player.Bonus = getBonuses("Заклинание Мгновенного Выздоровления", "instant_recovery_spell")
	player.Health = 1
	player.HealthMax = 25
	testCtx.DB.Save(player)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	reqBody := map[string]interface{}{
		"bonus":  "instant_recovery_spell",
		"option": "",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.Equal(t, player.HealthMax, player.Health)
}

func TestBonus_MagicDuckAntiMagic(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	battleStart := "player"
	section, enemies := helpers2.SetupBattleSection(t, testCtx.DB, 10, 2, &battleStart)
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

	player.Bonus = getBonuses("Магическая утка", "magic_duck")
	health := uint(30)
	player.Debuff = []entities.Debuff{
		{Alias: "magic", Health: &health},
		{Alias: "poison", Health: &health},
	}
	testCtx.DB.Save(player)

	assert.True(t, helpers.HasDebuff(player.Debuff, entities.AliasMagicReason))

	reqBody := map[string]interface{}{
		"bonus":  "magic_duck",
		"option": "anti_magic",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	resp = helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/get-section", token, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp = helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.False(t, helpers.HasDebuff(player.Debuff, entities.AliasMagicReason))

	enemy := helpers2.GetEnemyByEnemyID(t, testCtx.DB, enemies[0].ID)

	assert.True(t, helpers.HasDebuff(enemy.Debuff, entities.DebuffAliasMagicOffReason))
}

func TestBonus_MagicDuckSection(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	helpers2.CreatePlayerSection(t, testCtx.DB, entities.PlayerSection{
		SectionID: section.ID,
		PlayerID:  player.ID,
	})

	section158 := helpers2.CreateTestSection(t, testCtx.DB, 158, entities.SectionType("normal"))

	player.Bonus = getBonuses("Магическая утка", "magic_duck")
	testCtx.DB.Save(player)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	reqBody := map[string]interface{}{
		"bonus":  "magic_duck",
		"option": "section",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.Equal(t, section158.ID, player.SectionID)
}

func TestBonus_MagicRingLeft(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	helpers2.CreatePlayerSection(t, testCtx.DB, entities.PlayerSection{
		SectionID: section.ID,
		PlayerID:  player.ID,
	})

	player.Bonus = getBonuses("Магическое кольцо", "magic_ring")
	player.Health = 25
	player.HealthMax = 25
	testCtx.DB.Save(player)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	reqBody := map[string]interface{}{
		"bonus":  "magic_ring",
		"option": "left",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.Equal(t, player.HealthMax+25, player.Health)
}

func TestBonus_MagicRingRight(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	section := helpers2.CreateTestSection(t, testCtx.DB, 1, entities.SectionType("normal"))
	player := helpers2.CreateTestPlayer(t, testCtx.DB, user.ID, section.ID)
	helpers2.CreatePlayerSection(t, testCtx.DB, entities.PlayerSection{
		SectionID: section.ID,
		PlayerID:  player.ID,
	})

	player.Bonus = getBonuses("Магическое кольцо", "magic_ring")
	player.Weapons = []entities.Weapons{
		{Name: "Молнии", Damage: 10, MinCubeHit: 0, Item: "lightning", Count: 10},
	}
	testCtx.DB.Save(player)

	jwtSvc := helpers2.SetupTestAuthService(t)
	token := helpers2.CreateAuthenticatedRequest(t, user.ID, jwtSvc)

	reqBody := map[string]interface{}{
		"bonus":  "magic_ring",
		"option": "right",
	}

	resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
	assert.True(t, apiResp.Success)

	player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

	assert.Equal(t, player.Weapons[0].Count, uint(11))
}

func TestBonus_Wand(t *testing.T) {
	testCtx := helpers2.SetupRouteTest(t)
	defer helpers2.CleanupAllTables(testCtx.DB, testCtx.Context)

	user := helpers2.CreateTestUser(t, testCtx.DB)
	battleStart := "player"
	section, enemies := helpers2.SetupBattleSection(t, testCtx.DB, 10, 2, &battleStart)
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

	resultWithRandom := false
	for i := 0; i <= 100; i++ {
		player.Bonus = getBonuses("Волшебная палочка", "wand")
		testCtx.DB.Save(player)

		reqBody := map[string]interface{}{
			"bonus":  "wand",
			"option": "",
		}

		resp := helpers2.MakeRequest(t, testCtx.Server, "POST", "/api/game/ability/bonus", token, reqBody)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		apiResp := helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
		assert.True(t, apiResp.Success)

		player = helpers2.GetPlayerByPlayerID(t, testCtx.DB, player.ID)

		resp = helpers2.MakeRequest(t, testCtx.Server, "GET", "/api/game/get-section", token, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		apiResp = helpers2.ParseResponse[gameDto.CurrentResponse](t, resp)
		assert.True(t, apiResp.Success)

		enemy := helpers2.GetEnemyByEnemyID(t, testCtx.DB, enemies[0].ID)

		if helpers.HasDebuff(enemy.Debuff, entities.DebuffAliasSkipReason) {
			assert.Equal(t, *enemy.Debuff[0].Duration, uint(4))

			resultWithRandom = true

			break
		}
	}

	if !resultWithRandom {
		assert.Fail(t, "Магическая палочка не сработала")
	}
}

func getBonuses(bonusName string, bonusAlias string) []entities.PlayerBonus {
	return []entities.PlayerBonus{
		{Name: &bonusName, Alias: &bonusAlias},
	}
}
