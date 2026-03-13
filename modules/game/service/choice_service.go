package service

import (
	"context"
	"errors"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	gameDto "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	player2 "gamebook-backend/modules/game/player"
	"gamebook-backend/modules/game/repository"
	"gamebook-backend/modules/game/section"
	"gamebook-backend/modules/game/service/section/custom"

	"gorm.io/gorm"
)

type ChoiceService interface {
	Action(ctx context.Context, req gameDto.ActionRequest, player entities.Player) (gameDto.ActionResponse, error)
	Move(ctx context.Context, req gameDto.ActionRequest, player entities.Player) (gameDto.ActionResponse, error)
}

type choiceService struct {
	sectionRepository    repository.SectionRepository
	transitionRepository repository.TransitionRepository
	diceRepository       repository.DiceRepository
	playerRepository     repository.PlayerRepository
	db                   *gorm.DB
	playerUpdateListener listener.PlayerUpdateListener
}

func NewChoiceService(
	sectionRepo repository.SectionRepository,
	transitionRepo repository.TransitionRepository,
	diceRepo repository.DiceRepository,
	playerRepo repository.PlayerRepository,
	db *gorm.DB,
) ChoiceService {
	playerUpdateListener, _ := listener.HandleEvent(db, "player_update")

	return &choiceService{
		sectionRepository:    sectionRepo,
		transitionRepository: transitionRepo,
		diceRepository:       diceRepo,
		playerRepository:     playerRepo,
		db:                   db,
		playerUpdateListener: playerUpdateListener,
	}
}

func (s *choiceService) Action(ctx context.Context, req gameDto.ActionRequest, player entities.Player) (gameDto.ActionResponse, error) {
	if player.Health == 0 && player.Section.Number != gameDto.SectionDeath {
		deathSection, err := s.sectionRepository.GetBySectionNumber(ctx, s.db, gameDto.SectionDeath)
		if err != nil {
			return dto.ActionResponse{
				Result: gameDto.ResultFalse,
				Error:  fmt.Sprintf("%v", err),
			}, err
		}

		playerUpdate := player2.NewPlayerUpdate(s.db, player.ID)
		_, err = playerUpdate.Update(ctx, event.PlayerUpdateEvent{
			SectionID: &deathSection.ID,
		})
		if err != nil {
			return dto.ActionResponse{
				Result: gameDto.ResultFalse,
				Error:  fmt.Sprintf("%v", err),
			}, nil
		}

		return dto.ActionResponse{
			Result: gameDto.ResultTrue,
		}, nil
	}

	if custom.IsCustom(player.Section.Number) {
		requestSection, err := s.sectionRepository.GetByID(ctx, s.db, req.Transition)
		if err != nil {
			return dto.ActionResponse{
				Result: gameDto.ResultFalse,
				Error:  fmt.Sprintf("%v", err),
			}, err
		}

		if player.Section.Number == 156 || player.Section.Number == 158 {
			history := helper.GetVisitedSections(player.PlayerSection)

			for _, historyItem := range history {
				if requestSection.Number == historyItem {
					return dto.ActionResponse{
						Result:  gameDto.ResultTrue,
						Content: requestSection.Text,
						Actions: []gameDto.TransitionDTO{
							{Text: "Перейти в эту секцию", TransitionID: requestSection.ID},
						},
					}, err
				}
			}
		}

		if player.Section.Number == 157 || player.Section.Number == 158 {
			listSections := helper.GetNotVisitedSections(player.PlayerSection)

			for _, sectionItem := range listSections {
				if requestSection.Number == sectionItem {
					return dto.ActionResponse{
						Result:  gameDto.ResultTrue,
						Content: requestSection.Text,
						Actions: []gameDto.TransitionDTO{
							{Text: "Перейти в эту секцию", TransitionID: requestSection.ID},
						},
					}, err
				}
			}
		}
	}

	transition, err := s.transitionRepository.GetByTransitionID(ctx, s.db, req.Transition)
	if err != nil {
		return dto.ActionResponse{}, gameDto.ErrSectionNotFound
	}

	if transition.Section.Choice != nil {
		if *transition.Section.Choice.MaxSelections != uint(len(req.Data)) {
			return dto.ActionResponse{
				Result: gameDto.ResultFalse,
				Error:  fmt.Sprintf("Необходимо выбрать %d элементов", *transition.Section.Choice.MaxSelections),
			}, nil
		}

		bag := helper.GetFullBagItems(req.Data)
		playerUpdate := player2.NewPlayerUpdate(s.db, player.ID)
		playerChanged, err := playerUpdate.Update(ctx, event.PlayerUpdateEvent{
			PlayerID: player.ID,
			Bag:      &bag,
		})
		if err != nil {
			return dto.ActionResponse{
				Result: gameDto.ResultFalse,
				Error:  fmt.Sprintf("%v", err),
			}, err
		}

		player = *playerChanged
	}

	dices, err := s.diceRepository.GetLastByPlayerId(ctx, s.db, player.ID, gameDto.ReasonChoice)
	if err != nil && !errors.Is(err, gameDto.MessageDicesNotDefined) {
		return dto.ActionResponse{
			Result: gameDto.ResultFalse,
			Error:  fmt.Sprintf("%v", err),
		}, err
	}
	result, err := section.Check(ctx, s.db, transition, &dices, player)
	if err != nil {
		return dto.ActionResponse{
			Result: gameDto.ResultFalse,
			Error:  fmt.Sprintf("%v", err),
		}, err
	}

	if result {
		if transition.TargetSection.Number == 2 && player.Section.Number == 9 {
			playerRef, err := player2.ResetPlayer(ctx, s.db, player)
			if err != nil {
				return dto.ActionResponse{}, err

			}

			player = *playerRef
		}

		sectionUpdate := section.NewSectionUpdate(s.db)
		err := sectionUpdate.Update(ctx, player, transition)
		if err != nil {
			return gameDto.ActionResponse{}, err
		}
	}

	return dto.ActionResponse{
		Result: gameDto.ResultTrue,
	}, nil
}

func (s *choiceService) Move(ctx context.Context, req gameDto.ActionRequest, player entities.Player) (gameDto.ActionResponse, error) {
	if custom.IsCustom(player.Section.Number) {
		requestSection, err := s.sectionRepository.GetByID(ctx, s.db, req.Transition)
		if err != nil {
			return dto.ActionResponse{
				Result: gameDto.ResultFalse,
				Error:  fmt.Sprintf("%v", err),
			}, err
		}

		if player.Section.Number == 156 {
			history := helper.GetVisitedSections(player.PlayerSection)

			for _, historyItem := range history {
				if requestSection.Number == historyItem {
					err = s.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
						PlayerID:  player.ID,
						SectionID: &requestSection.ID,
					})
					if err != nil {
						return dto.ActionResponse{
							Result: gameDto.ResultFalse,
						}, err
					}

					return dto.ActionResponse{
						Result: gameDto.ResultTrue,
					}, nil
				}
			}
		}

		if player.Section.Number == 157 {
			listSections := helper.GetNotVisitedSections(player.PlayerSection)

			for _, sectionItem := range listSections {
				if requestSection.Number == sectionItem {
					err = s.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
						PlayerID:  player.ID,
						SectionID: &requestSection.ID,
					})
					if err != nil {
						return dto.ActionResponse{
							Result: gameDto.ResultFalse,
						}, err
					}

					return dto.ActionResponse{
						Result: gameDto.ResultTrue,
					}, nil
				}
			}
		}

		if player.Section.Number == 158 {
			listSections := helper.GetAllSections()

			for _, sectionItem := range listSections {
				if requestSection.Number == sectionItem {
					err = s.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
						PlayerID:  player.ID,
						SectionID: &requestSection.ID,
					})
					if err != nil {
						return dto.ActionResponse{
							Result: gameDto.ResultFalse,
						}, err
					}

					return dto.ActionResponse{
						Result: gameDto.ResultTrue,
					}, nil
				}
			}
		}
	}

	return dto.ActionResponse{
		Result: gameDto.ResultFalse,
		Error:  fmt.Sprintf("%v", gameDto.ErrCustomSectionNotFound),
	}, gameDto.ErrCustomSectionNotFound
}
