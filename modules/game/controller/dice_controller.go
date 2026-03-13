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
	DiceController interface {
		RollTheDice(ctx *gin.Context)
	}

	diceController struct {
		playerService service.PlayerService
		diceService   service.DiceService
	}
)

func NewDiceController(injector *do.Injector, ps service.PlayerService, ds service.DiceService) DiceController {
	return &diceController{
		playerService: ps,
		diceService:   ds,
	}
}

func (c *diceController) RollTheDice(ctx *gin.Context) {
	var err error

	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	diceDTO, err := c.diceService.RollTheDice(ctx, player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedRollTheDiceRequest, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessActionRequest, diceDTO)
	ctx.JSON(http.StatusOK, res)
}
