package bribe

import (
	"context"
	"gamebook-backend/database/entities"
	gameDTO "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func IsBribe(tx *gorm.DB, playerID uuid.UUID) (bool, error) {
	dicesRepository := repository.NewDiceRepository(tx)
	bribeService := NewBribeService(dicesRepository)

	return bribeService.IsBribe(context.Background(), tx, playerID)
}

func BribeAction(tx *gorm.DB, player entities.Player) error {
	dicesRepository := repository.NewDiceRepository(tx)
	bribeService := NewBribeService(dicesRepository)

	return bribeService.BribeAction(context.Background(), tx, player)
}

func GetBribeResultTransition(tx *gorm.DB, player entities.Player) (gameDTO.TransitionDTO, error) {
	dicesRepository := repository.NewDiceRepository(tx)
	bribeService := NewBribeService(dicesRepository)

	return bribeService.GetBribeResultTransition(context.Background(), tx, player)
}

func GetGotoSectionBribeSuccess(player entities.Player) (gameDTO.TransitionDTO, error) {
	for _, item := range player.Section.Transitions {
		if item.BribeResult == nil {
			continue
		}
		if *item.BribeResult {
			return gameDTO.TransitionDTO{
				Text:         item.Text,
				TransitionID: item.ID,
			}, nil
		}
	}

	return gameDTO.TransitionDTO{}, gameDTO.MessageBribeNotFoundSuccessTransition
}

func GetGotoSectionBribeFail(player entities.Player) (gameDTO.TransitionDTO, error) {
	for _, item := range player.Section.Transitions {
		if item.BribeResult == nil {
			continue
		}
		if !*item.BribeResult {
			return gameDTO.TransitionDTO{
				Text:         item.Text,
				TransitionID: item.ID,
			}, nil
		}
	}

	return gameDTO.TransitionDTO{}, gameDTO.MessageBribeNotFoundFailTransition
}
