package player

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"

	"gorm.io/gorm"
)

func ResetPlayer(ctx context.Context, db *gorm.DB, player entities.Player) (*entities.Player, error) {
	rollTheDices := dice.NewRollTheDices(db, &player)
	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, player)
	if err != nil {
		return nil, err
	}
	weapons := []entities.Weapons{
		{Name: "Руки", MinCubeHit: 6, Item: "hand"},
		{Name: "Меч-кладенец", Damage: 5, MinCubeHit: 4, Item: "sword"},
		{Name: "Молнии", Damage: 10, MinCubeHit: 0, Item: "lightning", Count: 10},
		{Name: "Шаровые молнии", Damage: 75, MinCubeHit: 6, Item: "ball lightning", Count: 2},
	}
	meds := entities.Meds{Name: "Лекарства", Item: "chain mail", Count: 15}

	newHealth := (*diceFirst + *diceSecond) * 4

	player.Weapons = weapons
	player.Meds = meds
	player.Bag = nil
	player.Health = newHealth
	player.HealthMax = newHealth
	player.Debuff = nil
	player.Buff = nil
	player.Bonus = helpBonuses()
	player.Gold = 0

	return &player, nil
}

func helpBonuses() []entities.PlayerBonus {
	bonusName1 := "Счастливый камушек"
	bonusAlias1 := "lucky_stone"

	return []entities.PlayerBonus{
		{Name: &bonusName1, Alias: &bonusAlias1},
	}
}
