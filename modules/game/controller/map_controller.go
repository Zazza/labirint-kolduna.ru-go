package controller

import (
	gameDto "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/service"
	"net/http"

	"gamebook-backend/pkg/constants"
	"gamebook-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type (
	MapController interface {
		GetMap(ctx *gin.Context)
	}

	mapController struct {
		playerService service.PlayerService
		mapService    service.MapService
		db            *gorm.DB
	}
)

func NewMapController(
	injector *do.Injector,
	playerService service.PlayerService,
	mapService service.MapService,
) MapController {
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	return &mapController{
		playerService: playerService,
		mapService:    mapService,
		db:            db,
	}
}

func (c *mapController) GetMap(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	player, err := c.playerService.GetByUserId(ctx, userId)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetPlayer, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if !helper.HasBagItem(player.Bag, "mapIngredients") {
		res := utils.BuildResponseFailed(gameDto.MessageMapIngredientNotFound, "", nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	result, err := c.mapService.GetMap(ctx.Request.Context(), c.db, player)
	if err != nil {
		res := utils.BuildResponseFailed(gameDto.MessageFailedGetMap, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(gameDto.MessageSuccessGetMap, result)
	ctx.JSON(http.StatusOK, res)
}
