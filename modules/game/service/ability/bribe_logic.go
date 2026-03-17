package ability

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bribe"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/template"
	"gorm.io/gorm"
)

type BribeLogic struct {
	db *gorm.DB
}

func NewBribeLogic(db *gorm.DB) *BribeLogic {
	return &BribeLogic{
		db: db,
	}
}

func (b *BribeLogic) Validate(ctx context.Context, player entities.Player) error {
	if player.Gold == 0 {
		return fmt.Errorf("player has no gold to bribe")
	}
	return nil
}

func (b *BribeLogic) Execute(ctx context.Context, player entities.Player, req dto.BribeRequest) error {
	rollTheDices := dice.NewRollTheDices(b.db, &player)
	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, player)
	if err != nil {
		return err
	}

	diceTemplate, err := template.GetDicesTemplate(ctx, *diceFirst, *diceSecond, true)
	if err != nil {
		return err
	}

	helper.SafeHTMLDescriptionMessageWithContext(ctx, player.ID, fmt.Sprintf("<p>%s</p>", diceTemplate))

	err = bribe.BribeAction(b.db, player)
	if err != nil {
		return err
	}

	err = rollTheDices.StoreDices(ctx, player, *diceFirst, *diceSecond, dto.ReasonBribe)
	if err != nil {
		return err
	}

	return nil
}

func (b *BribeLogic) Result() dto.BribeDTO {
	return dto.BribeDTO{
		Result:  true,
		Message: "Bribe completed successfully",
	}
}
