package providers

import (
	"context"
	"gamebook-backend/config"
	authController "gamebook-backend/modules/auth/controller"
	authRepo "gamebook-backend/modules/auth/repository"
	authService "gamebook-backend/modules/auth/service"
	"gamebook-backend/modules/game/channel"
	abilityController "gamebook-backend/modules/game/controller"
	battleController "gamebook-backend/modules/game/controller"
	choiceController "gamebook-backend/modules/game/controller"
	diceController "gamebook-backend/modules/game/controller"
	sectionController "gamebook-backend/modules/game/controller"
	playerLogListener "gamebook-backend/modules/game/listener"
	logService "gamebook-backend/modules/game/log"
	battleRepo "gamebook-backend/modules/game/repository"
	bonusRepo "gamebook-backend/modules/game/repository"
	diceRepo "gamebook-backend/modules/game/repository"
	playerRepo "gamebook-backend/modules/game/repository"
	playerSectionEnemyRepo "gamebook-backend/modules/game/repository"
	playerSectionRepo "gamebook-backend/modules/game/repository"
	sectionRepo "gamebook-backend/modules/game/repository"
	transitionRepo "gamebook-backend/modules/game/repository"
	abilityService "gamebook-backend/modules/game/service"
	choiceService "gamebook-backend/modules/game/service"
	diceService "gamebook-backend/modules/game/service"
	playerService "gamebook-backend/modules/game/service"
	"gamebook-backend/modules/game/service/battle"
	"gamebook-backend/modules/game/service/section"
	userController "gamebook-backend/modules/user/controller"
	"gamebook-backend/modules/user/repository"
	userService "gamebook-backend/modules/user/service"
	"gamebook-backend/pkg/constants"

	"github.com/samber/do"
	"gorm.io/gorm"
)

func InitDatabase(cfg *config.Config, injector *do.Injector) {
	do.ProvideNamed(injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
		return config.SetUpDatabaseConnection(cfg), nil
	})
}

func RegisterDependencies(cfg *config.Config, injector *do.Injector) {
	InitDatabase(cfg, injector)

	do.ProvideNamed(injector, constants.JWTService, func(i *do.Injector) (authService.JWTService, error) {
		return authService.NewJWTService(cfg), nil
	})

	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := authRepo.NewRefreshTokenRepository(db)

	sectionRepository := sectionRepo.NewSectionRepository(db)
	transitionRepository := transitionRepo.NewTransitionRepository(db)
	playerRepository := playerRepo.NewPlayerRepository(db)
	playerSectionRepository := playerSectionRepo.NewPlayerSectionRepository(db)
	playerSectionEnemyRepository := playerSectionEnemyRepo.NewPlayerSectionEnemyRepository(db)
	battleRepository := battleRepo.NewBattleRepository(db)
	diceRepository := diceRepo.NewDiceRepository(db)
	bonusRepository := bonusRepo.NewBonusRepository(db)
	playerLogRepository := playerRepo.NewPlayerLogRepository(db)

	eventChannel := channel.NewInMemoryEventChannel(nil)

	ctx := context.Background()

	playerLogListenerInstance := playerLogListener.NewPlayerLogListener(
		eventChannel,
		db,
		playerLogRepository,
	)
	go playerLogListenerInstance.Update(ctx)

	userServiceInstance := userService.NewUserService(userRepository, db)
	authServiceInstance := authService.NewAuthService(userRepository, refreshTokenRepository, jwtService, db)

	playerServiceInstance := playerService.NewPlayerService(playerRepository, playerSectionRepository, sectionRepository, db)

	playerLogServiceInstance := logService.NewPlayerLogService()

	bonusServiceInstance := abilityService.NewBonusService(bonusRepository)
	transitionServiceInstance := abilityService.NewTransitionService(transitionRepository, diceRepository, sectionRepository)
	playerProfileServiceInstance := abilityService.NewPlayerProfileService(bonusServiceInstance)
	bribeServiceInstance := abilityService.NewBribeService()
	diceServiceInstance := diceService.NewDiceServiceWithLogging(diceRepository, db, playerLogServiceInstance)

	sectionServiceInstance := section.NewSectionService(
		sectionRepository,
		transitionServiceInstance,
		bonusServiceInstance,
		playerProfileServiceInstance,
		diceServiceInstance,
		bribeServiceInstance,
	)
	battleSectionServiceInstance := section.NewBattleSectionService(
		sectionRepository,
		diceRepository,
		battleRepository,
		playerRepository,
		playerSectionEnemyRepository,
		db,
	)

	battleServiceInstance := battle.NewServiceWithLogging(sectionRepository, diceRepository, battleRepository, playerRepository, playerSectionEnemyRepository, db, playerLogServiceInstance)
	sleepySectionServiceInstance := section.NewSleepySectionService(
		sectionRepository,
		playerRepository,
		diceRepository,
		playerSectionRepository,
		db,
	)

	choiceServiceInstance := choiceService.NewChoiceService(sectionRepository, transitionRepository, diceRepository, playerRepository, db)

	mapServiceInstance := abilityService.NewMapService(sectionRepository, playerSectionRepository)

	abilityServiceInstance := abilityService.NewAbilityService(
		playerRepository,
		diceRepository,
		sectionRepository,
		battleRepository,
		transitionRepository,
		playerSectionRepository,
		playerSectionEnemyRepository,
		bonusRepository,
		db,
	)

	do.Provide(
		injector, func(i *do.Injector) (userController.UserController, error) {
			return userController.NewUserController(i, userServiceInstance), nil
		},
	)

	do.Provide(
		injector, func(i *do.Injector) (authController.AuthController, error) {
			return authController.NewAuthController(i, authServiceInstance), nil
		},
	)

	do.Provide(
		injector, func(i *do.Injector) (sectionController.SectionController, error) {
			return sectionController.NewSectionController(
				i,
				playerServiceInstance,
				battleSectionServiceInstance,
				sectionServiceInstance,
				sleepySectionServiceInstance,
			), nil
		},
	)

	do.Provide(
		injector, func(i *do.Injector) (diceController.DiceController, error) {
			return diceController.NewDiceController(i, playerServiceInstance, diceServiceInstance), nil
		},
	)

	do.Provide(
		injector, func(i *do.Injector) (battleController.BattleController, error) {
			return battleController.NewBattleController(i, playerServiceInstance, battleServiceInstance), nil
		},
	)

	do.Provide(
		injector, func(i *do.Injector) (choiceController.ChoiceController, error) {
			return choiceController.NewChoiceController(i, playerServiceInstance, choiceServiceInstance), nil
		},
	)

	do.Provide(
		injector, func(i *do.Injector) (abilityController.AbilityController, error) {
			return abilityController.NewAbilityController(i, playerServiceInstance, abilityServiceInstance), nil
		},
	)

	do.Provide(
		injector, func(i *do.Injector) (abilityController.MapController, error) {
			return abilityController.NewMapController(i, playerServiceInstance, mapServiceInstance), nil
		},
	)
}
