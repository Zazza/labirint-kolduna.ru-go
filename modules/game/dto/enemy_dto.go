package dto

const (
	BattleStartDice1 = "dice1"
	BattleStartDice2 = "dice2"
	BattleStartDice3 = "dice3"
	BattleStartDice4 = "dice4"
	BattleStartDice5 = "dice5"

	BattleStartDices  = "dices"
	BattleStartPlayer = "player"
	BattleStartEnemy  = "enemy"
)

type EnemyDamageDTO struct {
	Damage      uint
	Description string
}
