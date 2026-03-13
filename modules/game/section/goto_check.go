package section

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/expression"
	"strings"

	"gorm.io/gorm"
)

func Check(ctx context.Context, db *gorm.DB, item entities.Transition, dices *entities.Dice, player entities.Player) (bool, error) {
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
	if dices == nil {
		return true, nil
	} else if item.Dice != nil {
		dicesResult = dices.DiceFirst
		itemExpression = *item.Dice
	} else if item.Dices != nil {
		dicesResult = dices.DiceFirst + dices.DiceSecond
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
		bagItem := CheckConditions(ctx, db, item.Condition, &player)
		outputWithCondition, err = expression.Run(fmt.Sprintf("%t && %t", output, bagItem))

		output = outputWithCondition
	}

	if item.PlayerChange != nil && item.PlayerChange.ReturnToSection != nil {
		outputWithReturnToSection := *item.PlayerChange.ReturnToSection == player.ReturnToSection
		output = outputWithReturnToSection && output.(bool)
	}

	return output.(bool), nil
}

func CheckSimple(ctx context.Context, db *gorm.DB, item entities.Transition, player entities.Player) (bool, error) {
	if item.AvailableOnce {
		for _, logSection := range player.PlayerSection {
			if logSection.SectionID == item.TargetSectionID {
				return false, nil
			}
		}
	}

	output := true

	if item.Condition != nil {
		conditionResult := CheckConditions(ctx, db, item.Condition, &player)
		output = conditionResult
	}

	if item.PlayerChange != nil && item.PlayerChange.ReturnToSection != nil {
		outputWithReturnToSection := *item.PlayerChange.ReturnToSection == player.ReturnToSection
		output = outputWithReturnToSection && output
	}

	return output, nil
}
