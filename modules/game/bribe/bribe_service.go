package bribe

import (
	"context"
	"errors"
	"fmt"
	"gamebook-backend/database/entities"
	dice2 "gamebook-backend/modules/game/dice"
	gameDTO "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/expression"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/log"
	"gamebook-backend/modules/game/repository"
	template2 "gamebook-backend/modules/game/template"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BribeService interface {
	IsBribe(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) (bool, error)
	BribeAction(ctx context.Context, tx *gorm.DB, player entities.Player) error
	GetBribeResultTransition(ctx context.Context, tx *gorm.DB, player entities.Player) (gameDTO.TransitionDTO, error)
}

type bribeService struct {
	diceRepository repository.DiceRepository
	logService     log.PlayerLogService
}

func NewBribeService(
	diceRepo repository.DiceRepository,
) BribeService {
	return &bribeService{
		diceRepository: diceRepo,
	}
}

func NewBribeServiceWithLogging(
	diceRepo repository.DiceRepository,
	logService log.PlayerLogService,
) BribeService {
	return &bribeService{
		diceRepository: diceRepo,
		logService:     logService,
	}
}

func (s *bribeService) IsBribe(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) (bool, error) {
	_, err := s.diceRepository.GetLastByPlayerId(ctx, tx, playerID, gameDTO.ReasonBribe)
	if errors.Is(err, gameDTO.MessageDicesNotDefined) {
		return false, nil
	} else if err != nil && !errors.Is(err, gameDTO.MessageDicesNotDefined) {
		return false, err
	}

	return true, err
}

func (s *bribeService) BribeAction(ctx context.Context, tx *gorm.DB, player entities.Player) error {
	playerUpdateListener, err := listener.HandleEvent(tx, "player_update")
	if err != nil {
		return err
	}

	if player.Section.Bribe.Amount != nil {
		success := *player.Section.Bribe.Amount <= player.Gold
		if success {
			player.Gold -= *player.Section.Bribe.Amount
			err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
				PlayerID: player.ID,
				Gold:     &player.Gold,
			})
			if err != nil {
				return err
			}
		}

		if s.logService != nil {
			s.logService.LogBribe(player.ID, success, *player.Section.Bribe.Amount, player.Section.ID.String())
		}

		return nil
	}

	if player.Section.Bribe.AmountDice != nil {
		rollTheDices := dice2.NewRollTheDices(tx, &player)
		diceFirst, err := rollTheDices.RollTheDice(context.Background(), player)
		if err != nil {
			return err
		}

		template, err := template2.GetDiceTemplate(context.Background(), *diceFirst, true)
		if err != nil {
			return err
		}

		helper.SafeHTMLDescriptionMessage(player.ID, fmt.Sprintf("<p>Размер взятки: %s</p>", template))

		success := *diceFirst <= player.Gold
		if success {
			player.Gold -= *diceFirst
			err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
				PlayerID: player.ID,
				Gold:     &player.Gold,
			})
			if err != nil {
				return err
			}
		}

		if s.logService != nil {
			s.logService.LogBribe(player.ID, success, *diceFirst, player.Section.ID.String())
		}

		return nil
	}

	if player.Section.Bribe.AmountDices != nil {
		rollTheDices := dice2.NewRollTheDices(tx, &player)
		diceFirst, diceSecond, err := rollTheDices.RollTheDices(context.Background(), player)
		if err != nil {
			return err
		}

		template, err := template2.GetDicesTemplate(context.Background(), *diceFirst, *diceSecond, true)
		if err != nil {
			return err
		}

		helper.SafeHTMLDescriptionMessage(player.ID, fmt.Sprintf("<p>Размер взятки: %s</p>", template))

		success := *diceFirst+*diceSecond <= player.Gold
		if success {
			player.Gold -= *diceFirst + *diceSecond
			err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
				PlayerID: player.ID,
				Gold:     &player.Gold,
			})
			if err != nil {
				return err
			}
		}

		if s.logService != nil {
			s.logService.LogBribe(player.ID, success, *diceFirst+*diceSecond, player.Section.ID.String())
		}

		return nil
	}

	return gameDTO.MessageBribeLogicError
}

func (s *bribeService) GetBribeResultTransition(ctx context.Context, tx *gorm.DB, player entities.Player) (gameDTO.TransitionDTO, error) {
	dice, err := s.diceRepository.GetLastByPlayerId(ctx, tx, player.ID, gameDTO.ReasonBribe)
	if err != nil {
		return gameDTO.TransitionDTO{}, err
	}

	if player.Section.Bribe.MinDiceHit != nil {
		result, err := expression.RunAndReturnBoolean(
			fmt.Sprintf("%d %s", dice.DiceFirst, *player.Section.Bribe.MinDiceHit),
		)
		if err != nil {
			return gameDTO.TransitionDTO{}, err
		}
		if !result {
			gotoSection, err := GetGotoSectionBribeFail(player)
			if err != nil {
				return gameDTO.TransitionDTO{}, err
			}

			return gotoSection, nil
		}
	}

	if player.Section.Bribe.MinDicesHit != nil {
		result, err := expression.RunAndReturnBoolean(
			fmt.Sprintf("%d %s", dice.DiceFirst+dice.DiceSecond, *player.Section.Bribe.MinDicesHit),
		)
		if err != nil {
			return gameDTO.TransitionDTO{}, err
		}
		if !result {
			gotoSection, err := GetGotoSectionBribeFail(player)
			if err != nil {
				return gameDTO.TransitionDTO{}, err
			}

			return gotoSection, nil
		}
	}

	if player.Section.Bribe.Amount != nil {
		if *player.Section.Bribe.Amount <= player.Gold {
			result, err := success(player, *player.Section.Bribe.Amount)
			if err != nil {
				return gameDTO.TransitionDTO{}, err
			}
			return result, nil
		}

		result, err := fail(player)
		if err != nil {
			return gameDTO.TransitionDTO{}, err
		}
		return result, nil
	}

	if player.Section.Bribe.AmountDice != nil {
		rollTheDices := dice2.NewRollTheDices(tx, &player)
		diceFirst, err := rollTheDices.RollTheDice(context.Background(), player)
		if err != nil {
			return gameDTO.TransitionDTO{}, err
		}

		if *diceFirst <= player.Gold {
			result, err := success(player, *diceFirst)
			if err != nil {
				return gameDTO.TransitionDTO{}, err
			}
			return result, nil
		}

		result, err := fail(player)
		if err != nil {
			return gameDTO.TransitionDTO{}, err
		}
		return result, nil
	}

	if player.Section.Bribe.AmountDices != nil {
		rollTheDices := dice2.NewRollTheDices(tx, &player)
		diceFirst, diceSecond, err := rollTheDices.RollTheDices(context.Background(), player)
		if err != nil {
			return gameDTO.TransitionDTO{}, err
		}

		if *diceFirst+*diceSecond <= player.Gold {
			result, err := success(player, *diceFirst+*diceSecond)
			if err != nil {
				return gameDTO.TransitionDTO{}, err
			}
			return result, nil
		}

		result, err := fail(player)
		if err != nil {
			return gameDTO.TransitionDTO{}, err
		}
		return result, nil
	}

	return gameDTO.TransitionDTO{}, gameDTO.MessageBribeLogicError
}

func success(player entities.Player, gold uint) (gameDTO.TransitionDTO, error) {
	gotoSection, err := GetGotoSectionBribeSuccess(player)
	if err != nil {
		return gameDTO.TransitionDTO{}, err
	}

	return gotoSection, nil
}

func fail(player entities.Player) (gameDTO.TransitionDTO, error) {
	gotoSection, err := GetGotoSectionBribeFail(player)
	if err != nil {
		return gameDTO.TransitionDTO{}, err
	}

	return gotoSection, nil
}
