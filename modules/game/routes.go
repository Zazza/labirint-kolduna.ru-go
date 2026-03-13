package game

import (
	authService "gamebook-backend/modules/auth/service"
	"gamebook-backend/modules/game/controller"
	"gamebook-backend/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"

	"gamebook-backend/middlewares"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	abilityController := do.MustInvoke[controller.AbilityController](injector)
	sectionController := do.MustInvoke[controller.SectionController](injector)
	battleController := do.MustInvoke[controller.BattleController](injector)
	choiceController := do.MustInvoke[controller.ChoiceController](injector)
	diceController := do.MustInvoke[controller.DiceController](injector)
	mapController := do.MustInvoke[controller.MapController](injector)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	gameRoutes := server.Group("/api/game")
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
			"/move",
			middlewares.Authenticate(jwtService),
			choiceController.Move,
		)
		gameRoutes.POST(
			"/battle",
			middlewares.Authenticate(jwtService),
			battleController.Battle,
		)
		gameRoutes.GET(
			"/profile",
			middlewares.Authenticate(jwtService),
			sectionController.GetProfile,
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
		gameRoutes.POST(
			"/ability/sleep/choice",
			middlewares.Authenticate(jwtService),
			abilityController.SleepChoice,
		)
		gameRoutes.POST(
			"/ability/bribe",
			middlewares.Authenticate(jwtService),
			abilityController.Bribe,
		)
		gameRoutes.GET(
			"/map",
			middlewares.Authenticate(jwtService),
			mapController.GetMap,
		)
	}
}
