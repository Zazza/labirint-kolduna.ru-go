package controller

import (
	gameDto "gamebook-backend/modules/game/dto"
	"net/http"

	"gamebook-backend/modules/game/service"
	"gamebook-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type (
	AbilityController interface {
		Meds(ctx *gin.Context)
		Bonus(ctx *gin.Context)
		Sleep(ctx *gin.Context)
		SleepChoice(ctx *gin.Context)
		Bribe(ctx *gin.Context)
	}

	abilityController struct {
		playerService  service.PlayerService
		abilityService service.AbilityService
	}
)

func NewAbilityController(
	injector *do.Injector,
	ps service.PlayerService,
	as service.AbilityService,
) AbilityController {
	return &abilityController{
		playerService:  ps,
		abilityService: as,
	}
}

func (c *abilityController) Meds(ctx *gin.Context) {
	var err error

	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.abilityService.Meds(ctx.Request.Context(), player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedCurrentRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessCurrentRequest, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *abilityController) Bonus(ctx *gin.Context) {
	var err error

	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req gameDto.BonusRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedActionRequest, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.abilityService.Bonus(ctx.Request.Context(), req, player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedCurrentRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessCurrentRequest, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *abilityController) Sleep(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.abilityService.Sleep(ctx.Request.Context(), player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedCurrentRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessCurrentRequest, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *abilityController) SleepChoice(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.abilityService.SleepChoice(ctx.Request.Context(), player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedCurrentRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessCurrentRequest, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *abilityController) Bribe(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.abilityService.Bribe(ctx.Request.Context(), player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedCurrentRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessCurrentRequest, result)
	ctx.JSON(http.StatusOK, res)
}
