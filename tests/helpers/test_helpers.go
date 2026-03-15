package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"gamebook-backend/config"
	"gamebook-backend/database/entities"
	"gamebook-backend/middlewares"
	authService "gamebook-backend/modules/auth/service"
	"gamebook-backend/modules/game/controller"
	gameDto "gamebook-backend/modules/game/dto"
	"gamebook-backend/pkg/constants"
	"gamebook-backend/providers"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestDB struct {
	DB *gorm.DB
}

type TestServer struct {
	Server *httptest.Server
}

type TestContext struct {
	Injector *do.Injector
	DB       *gorm.DB
	Server   *TestServer
	Context  context.Context
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewTestDB(injector *do.Injector) *TestDB {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	return &TestDB{DB: db}
}

func NewTestServer(injector *do.Injector) *TestServer {
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	router := setupGinTestServer(injector, jwtService)

	server := httptest.NewServer(router)

	return &TestServer{Server: server}
}

func setupGinTestServer(injector *do.Injector, jwtService authService.JWTService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(middlewares.ErrorLoggerMiddleware())
	router.Use(middlewares.CORSMiddleware())

	sectionController := do.MustInvoke[controller.SectionController](injector)
	battleController := do.MustInvoke[controller.BattleController](injector)
	choiceController := do.MustInvoke[controller.ChoiceController](injector)
	diceController := do.MustInvoke[controller.DiceController](injector)
	abilityController := do.MustInvoke[controller.AbilityController](injector)

	gameRoutes := router.Group("/api/game")
	{
		gameRoutes.GET(
			"/get-section",
			middlewares.Authenticate(jwtService),
			sectionController.GetSection,
		)
		gameRoutes.GET(
			"/role-the-dice",
			middlewares.Authenticate(jwtService),
			diceController.RollTheDice,
		)
		gameRoutes.POST(
			"/choice",
			middlewares.Authenticate(jwtService),
			choiceController.Action,
		)
		gameRoutes.POST(
			"/battle",
			middlewares.Authenticate(jwtService),
			battleController.Battle,
		)
		gameRoutes.POST(
			"/ability/meds",
			middlewares.Authenticate(jwtService),
			abilityController.Meds,
		)
		gameRoutes.POST(
			"/ability/bonus",
			middlewares.Authenticate(jwtService),
			abilityController.Bonus,
		)
		gameRoutes.POST(
			"/ability/sleep",
			middlewares.Authenticate(jwtService),
			abilityController.Sleep,
		)
	}

	return router
}

func CreateTestUser(t *testing.T, db *gorm.DB) *entities.User {
	userID := uuid.New()
	user := &entities.User{
		ID:       userID,
		Name:     "Test User",
		Password: "password123",
	}
	err := db.Create(user).Error
	require.NoError(t, err)
	return user
}

func CreateTestPlayer(t *testing.T, db *gorm.DB, userID uuid.UUID, sectionID uuid.UUID) *entities.Player {
	playerID := uuid.New()
	player := &entities.Player{
		ID:        playerID,
		UserID:    userID,
		SectionID: sectionID,
		Health:    20,
		HealthMax: 20,
		Gold:      50,
		Weapons: []entities.Weapons{
			{Name: "Sword", Damage: 10, MinCubeHit: 4, Item: "sword"},
		},
		Meds: entities.Meds{
			Name:  "Potion",
			Item:  "potion",
			Count: 2,
		},
		Bag: []entities.Bag{
			{Name: "Torch", Description: "A wooden torch"},
		},
	}
	err := db.Create(player).Error
	require.NoError(t, err)
	return player
}

func GetPlayerByPlayerID(t *testing.T, db *gorm.DB, PlayerID uuid.UUID) *entities.Player {
	var player entities.Player
	err := db.Where("id = ?", PlayerID).First(&player).Error
	require.NoError(t, err, "Player with ID %d not found", PlayerID)
	return &player
}

func GetEnemyByEnemyID(t *testing.T, db *gorm.DB, EnemyID uuid.UUID) *entities.PlayerSectionEnemy {
	var enemy entities.PlayerSectionEnemy
	err := db.Where("enemy_id = ?", EnemyID).First(&enemy).Error
	require.NoError(t, err, "Enemy with ID %d not found", EnemyID)
	return &enemy
}

func CreateTestEnemy(t *testing.T, db *gorm.DB, alias string) *entities.Enemy {
	enemyID := uuid.New()
	enemy := &entities.Enemy{
		ID:          enemyID,
		Alias:       alias,
		Name:        "Test Enemy",
		Damage:      5,
		MinDiceHits: 6,
		Health:      10,
		Defence:     2,
		PlayerArmor: true,
	}
	err := db.Create(enemy).Error
	require.NoError(t, err)
	return enemy
}

func AddEnemiesToSection(t *testing.T, db *gorm.DB, section *entities.Section, enemies []*entities.Enemy) {
	for _, enemy := range enemies {
		err := db.Model(section).Association("SectionEnemies").Append(enemy)
		require.NoError(t, err)
	}
}

func CreateTestTransition(t *testing.T, db *gorm.DB, sectionID, targetSectionID uuid.UUID, text string) *entities.Transition {
	transitionID := uuid.New()
	transition := &entities.Transition{
		ID:              transitionID,
		TextOrder:       transitionOrder,
		SectionID:       sectionID,
		TargetSectionID: targetSectionID,
		Text:            text,
		AvailableOnce:   false,
	}
	err := db.Create(transition).Error
	require.NoError(t, err)
	transitionOrder++
	return transition
}

func SetupBattleSection(t *testing.T, db *gorm.DB, number uint, enemyCount uint, battleStart *string) (*entities.Section, []*entities.Enemy) {
	sectionID := uuid.New()

	var battleSteps []*string

	enemies := make([]*entities.Enemy, enemyCount)
	for i := 0; uint(i) < enemyCount; i++ {
		enemies[i] = CreateTestEnemy(t, db, "enemy_"+uuid.New().String()[:8])

		battleSteps = []*string{StringPtr("player"), StringPtr(enemies[i].Alias)}
	}

	section := &entities.Section{
		ID:          sectionID,
		Type:        gameDto.SectionTypeBattle,
		Number:      number,
		Text:        "Battle section",
		BattleStart: battleStart,
		BattleSteps: battleSteps,
		Dices:       false,
	}
	err := db.Create(section).Error
	require.NoError(t, err)

	AddEnemiesToSection(t, db, section, enemies)

	return section, enemies
}

func MakeRequest(t *testing.T, server *TestServer, method, path string, token string, body interface{}) *http.Response {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, server.Server.URL+path, reqBody)
	require.NoError(t, err)

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func ParseResponse[T any](t *testing.T, resp *http.Response) APIResponse {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	require.NoError(t, err)

	apiResp.Success = resp.StatusCode == http.StatusOK

	return apiResp
}

func ParseData[T any](t *testing.T, resp *http.Response) T {
	apiResp := ParseResponse[T](t, resp)

	jsonData, err := json.Marshal(apiResp.Data)
	require.NoError(t, err)

	var result T
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	return result
}

var (
	enemyCounter    uint = 0
	transitionOrder uint = 1
)

func SetupRouteTest(t *testing.T) *TestContext {
	cfg, err := LoadConfig("../..")
	if err != nil {
		log.Fatalf("❌ Configuration error: %v", err)
	}

	injector := do.New()
	providers.RegisterDependencies(cfg, injector)

	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	ctx := context.Background()

	CleanupAllTables(db, ctx)

	transitionOrder = 1

	return &TestContext{
		Injector: injector,
		DB:       db,
		Server:   NewTestServer(injector),
		Context:  ctx,
	}
}

func LoadConfig(rootPath string) (*config.Config, error) {
	err := os.Setenv("APP_ENV", "test")
	if err != nil {
		return nil, err
	}

	config.LoadEnv(rootPath)
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func UIntPtr(u uint) *uint {
	return &u
}

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func BuffOrDebuffAliasPtr(alias entities.BuffOrDebuffAlias) *entities.BuffOrDebuffAlias {
	return &alias
}

func CreateTestSection(t *testing.T, db *gorm.DB, number uint, sectionType entities.SectionType) *entities.Section {
	sectionID := uuid.New()
	section := &entities.Section{
		ID:          sectionID,
		Type:        sectionType,
		Number:      number,
		Text:        "Test section text",
		Dices:       false,
		BattleStart: nil,
	}
	err := db.Create(section).Error
	require.NoError(t, err)
	return section
}
