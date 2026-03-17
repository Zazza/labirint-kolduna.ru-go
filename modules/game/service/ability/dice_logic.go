package ability

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"
	"gorm.io/gorm"
)

type DiceLogic struct {
	db             *gorm.DB
	diceRepository repository.DiceRepository
}

func NewDiceLogic(db *gorm.DB, diceRepo repository.DiceRepository) *DiceLogic {
	return &DiceLogic{
		db:             db,
		diceRepository: diceRepo,
	}
}

func (d *DiceLogic) Validate(ctx context.Context, player entities.Player) error {
	return nil
}

func (d *DiceLogic) Execute(ctx context.Context, player entities.Player) (dto.DiceDTO, error) {
	battleDicesDTO := d.diceRepository.FindBattleDicesByPlayerId(ctx, d.db, player.ID)
	if battleDicesDTO.Error != nil {
		return dto.DiceDTO{}, battleDicesDTO.Error
	}

	if len(player.Section.SectionEnemies) > 0 && battleDicesDTO.Exists && battleDicesDTO.Dices.DiceFirst != battleDicesDTO.Dices.DiceSecond {
		return dto.DiceDTO{}, dto.MessageBattleDicesAlreadyExist
	}

	rollTheDices := dice.NewRollTheDices(d.db, &player)
	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, player)
	if err != nil {
		return dto.DiceDTO{}, err
	}

	reason := dto.ReasonChoice
	if len(player.Section.SectionEnemies) > 0 {
		reason = dto.ReasonBattle
	}

	err = rollTheDices.StoreDices(ctx, player, *diceFirst, *diceSecond, reason)
	if err != nil {
		return dto.DiceDTO{}, err
	}

	return dto.DiceDTO{
		DiceFirst:  *diceFirst,
		DiceSecond: *diceSecond,
		Result:     dto.ResultTrue,
	}, nil
}

func (d *DiceLogic) Result() dto.DiceResultDTO {
	return dto.DiceResultDTO{
		DiceFirst:  0,
		DiceSecond: 0,
		Result:     true,
	}
}
