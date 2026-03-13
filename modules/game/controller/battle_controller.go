package controller

import (
	gameDto "gamebook-backend/modules/game/dto"
	"net/http"

	"gamebook-backend/modules/game/service"
	"gamebook-backend/modules/game/service/battle"
	"gamebook-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type (
	BattleController interface {
		Battle(ctx *gin.Context)
	}

	battleController struct {
		playerService service.PlayerService
		battleService battle.Service
	}
)

func NewBattleController(injector *do.Injector, ps service.PlayerService, bs battle.Service) BattleController {
	return &battleController{
		playerService: ps,
		battleService: bs,
	}
}

func (c *battleController) Battle(ctx *gin.Context) {
	var err error

	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req gameDto.ActionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedActionRequest, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.battleService.Action(ctx.Request.Context(), player, req)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedActionRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessActionRequest, result)
	ctx.JSON(http.StatusOK, res)
}
