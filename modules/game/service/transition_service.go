package service

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/expression"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/repository"
	sectionpkg "gamebook-backend/modules/game/section"
	template2 "gamebook-backend/modules/game/template"
	"strings"

	"gorm.io/gorm"
)

type TransitionService interface {
	Check(ctx context.Context, db *gorm.DB, transition *entities.Transition, player *entities.Player) (bool, error)
	CheckSimple(ctx context.Context, db *gorm.DB, transition *entities.Transition, player *entities.Player) (bool, error)
	GetAvailableTransitions(ctx context.Context, db *gorm.DB, player *entities.Player, dices *entities.Dice) ([]dto.TransitionDTO, error)
}

type transitionService struct {
	transitionRepository repository.TransitionRepository
	diceRepository       repository.DiceRepository
	sectionRepository    repository.SectionRepository
}

func NewTransitionService(
	transitionRepo repository.TransitionRepository,
	diceRepo repository.DiceRepository,
	sectionRepo repository.SectionRepository,
) TransitionService {
	return &transitionService{
		transitionRepository: transitionRepo,
		diceRepository:       diceRepo,
		sectionRepository:    sectionRepo,
	}
}

func (s *transitionService) Check(ctx context.Context, db *gorm.DB, item *entities.Transition, player *entities.Player) (bool, error) {
	var output, outputWithCondition any
	var err error

	if item.AvailableOnce {
		for _, logSection := range player.PlayerSection {
			if logSection.SectionID == item.TargetSectionID {
				return false, nil
			}
		}
	}

	var dicesResult uint
	var itemExpression []string

	playerDices, ok := ctx.Value("dices").(*entities.Dice)
	if !ok || playerDices == nil {
		return true, nil
	}

	if item.Dice != nil {
		dicesResult = playerDices.DiceFirst
		itemExpression = *item.Dice
	} else if item.Dices != nil {
		dicesResult = playerDices.DiceFirst + playerDices.DiceSecond
		itemExpression = *item.Dices
	} else {
		return true, nil
	}

	var expressionText []string
	for _, expressionString := range itemExpression {
		expressionText = append(expressionText, fmt.Sprintf("%d %s", dicesResult, expressionString))
	}
	output, err = expression.Run(strings.Join(expressionText, " && "))
	if err != nil {
		return false, err
	}

	if !output.(bool) {
		return false, nil
	}

	if item.Condition != nil {
		bagItem := sectionpkg.CheckConditions(ctx, db, item.Condition, player)
		outputWithCondition, err = expression.Run(fmt.Sprintf("%t && %t", output, bagItem))

		output = outputWithCondition
	}

	if item.PlayerChange != nil && item.PlayerChange.ReturnToSection != nil {
		outputWithReturnToSection := *item.PlayerChange.ReturnToSection == player.ReturnToSection
		output = outputWithReturnToSection && output.(bool)
	}

	return output.(bool), nil
}

func (s *transitionService) CheckSimple(ctx context.Context, db *gorm.DB, item *entities.Transition, player *entities.Player) (bool, error) {
	if item.AvailableOnce {
		for _, logSection := range player.PlayerSection {
			if logSection.SectionID == item.TargetSectionID {
				return false, nil
			}
		}
	}

	output := true

	if item.Condition != nil {
		conditionResult := sectionpkg.CheckConditions(ctx, db, item.Condition, player)
		output = conditionResult
	}

	if item.PlayerChange != nil && item.PlayerChange.ReturnToSection != nil {
		outputWithReturnToSection := *item.PlayerChange.ReturnToSection == player.ReturnToSection
		output = outputWithReturnToSection && output
	}

	return output, nil
}

func (s *transitionService) GetAvailableTransitions(ctx context.Context, db *gorm.DB, player *entities.Player, dices *entities.Dice) ([]dto.TransitionDTO, error) {
	var gotoItems []dto.TransitionDTO
	var ctxWithDices context.Context

	if dices != nil {
		ctxWithDices = context.WithValue(ctx, "dices", dices)
	} else {
		ctxWithDices = ctx
	}

	for _, item := range player.Section.Transitions {
		var result bool
		var err error

		if dices != nil {
			result, err = s.Check(ctxWithDices, db, &item, player)
		} else {
			result, err = s.CheckSimple(ctx, db, &item, player)
		}

		if err != nil {
			return nil, err
		}

		if result {
			sectionEntity, err := s.sectionRepository.GetByID(ctx, db, item.TargetSectionID)
			if err != nil {
				return nil, err
			}

			gotoItems = append(gotoItems, dto.TransitionDTO{
				Text:         item.Text,
				TransitionID: item.ID,
				Visited:      helper.IsVisited(sectionEntity.Number, player.PlayerSection),
			})
		}
	}

	return gotoItems, nil
}

func (s *transitionService) GetTransitionTemplate(ctx context.Context, transition *entities.Transition, dices *entities.Dice, player *entities.Player) (string, error) {
	if transition.Dice != nil {
		template, err := template2.GetDiceTemplate(
			ctx,
			dices.DiceFirst,
			true,
		)
		if err != nil {
			return "", err
		}
		return template, nil
	} else if transition.Dices != nil {
		template, err := template2.GetDicesTemplate(
			ctx,
			dices.DiceFirst,
			dices.DiceSecond,
			true,
		)
		if err != nil {
			return "", err
		}
		return template, nil
	}
	return "", nil
}
