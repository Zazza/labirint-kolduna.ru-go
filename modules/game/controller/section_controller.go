package controller

import (
	gameDto "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/service/section"
	"net/http"

	"gamebook-backend/modules/game/service"
	"gamebook-backend/modules/game/validation"
	"gamebook-backend/pkg/constants"
	"gamebook-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type (
	SectionController interface {
		GetSection(ctx *gin.Context)
		GetProfile(ctx *gin.Context)
	}

	sectionController struct {
		playerService        service.PlayerService
		battleSectionService section.BattleSectionService
		gameService          section.SectionService
		sleepyService        section.SleepySectionService
		gameValidation       *validation.GameValidation
		db                   *gorm.DB
	}
)

func NewSectionController(
	injector *do.Injector,
	ps service.PlayerService,
	bs section.BattleSectionService,
	gs section.SectionService,
	ss section.SleepySectionService,
) SectionController {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	gameValidation := validation.NewAuthValidation()
	return &sectionController{
		playerService:        ps,
		battleSectionService: bs,
		gameService:          gs,
		sleepyService:        ss,
		gameValidation:       gameValidation,
		db:                   db,
	}
}

func (c *sectionController) GetSection(ctx *gin.Context) {
	var err error

	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var result gameDto.CurrentResponse

	if len(player.Section.SectionEnemies) > 0 {
		result, err = c.battleSectionService.GetActivityByPlayer(ctx.Request.Context(), player)
		if err != nil {
			res := utils.BuildResponseFailed(gameDto.MessageFailedCurrentRequest, err.Error(), nil)
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
	} else if player.Section.Type == gameDto.SectionTypeSleepy {
		result, err = c.sleepyService.GetSection(ctx.Request.Context(), player)
		if err != nil {
			res := utils.BuildResponseFailed(gameDto.MessageFailedCurrentRequest, err.Error(), nil)
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
	} else {
		result, err = c.gameService.GetSection(ctx.Request.Context(), c.db, player)
		if err != nil {
			res := utils.BuildResponseFailed(gameDto.MessageFailedCurrentRequest, err.Error(), nil)
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
	}

	if result.Player.Health == 0 && player.Section.Number != gameDto.SectionDeath {
		result.Transitions = []gameDto.TransitionDTO{
			{
				Text:  "Ты погиб... Секция 9",
				Death: true,
			},
		}
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessCurrentRequest, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *sectionController) GetProfile(ctx *gin.Context) {
	var err error

	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.gameService.GetProfile(ctx.Request.Context(), c.db, player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedCurrentRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessCurrentRequest, result)
	ctx.JSON(http.StatusOK, res)
}
