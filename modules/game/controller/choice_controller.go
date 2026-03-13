package controller

import (
	gameDto "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/service"
	"gamebook-backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type (
	ChoiceController interface {
		Action(ctx *gin.Context)
		Move(ctx *gin.Context)
	}

	choiceController struct {
		playerService service.PlayerService
		choiceService service.ChoiceService
	}
)

func NewChoiceController(injector *do.Injector, ps service.PlayerService, bs service.ChoiceService) ChoiceController {
	return &choiceController{
		playerService: ps,
		choiceService: bs,
	}
}

func (c *choiceController) Action(ctx *gin.Context) {
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

	result, err := c.choiceService.Action(ctx.Request.Context(), req, player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedActionRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessActionRequest, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *choiceController) Move(ctx *gin.Context) {
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

	result, err := c.choiceService.Move(ctx.Request.Context(), req, player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedActionRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessActionRequest, result)
	ctx.JSON(http.StatusOK, res)
}
