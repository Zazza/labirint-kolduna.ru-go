package battle

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bonus/helpers"
	battleDTO "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/expression"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
)

type Player interface {
	GetDamageByWeapon(
		weapon string,
		player entities.Player,
		enemyDTO *battleDTO.EnemyDTO,
		gameCubeResultFirst uint,
		gameCubeResultSecond uint,
	) (uint, error)
	GetPlayerHealth(playerHealth uint, damage uint) uint
	GetPlayerWeapon(weapon string) []entities.Weapons
	UpdatePlayerAfterStep() error
}

type player struct {
	player               *entities.Player
	lastBattleLog        *entities.Battle
	battleLog            *[]entities.Battle
	dices                *entities.Dice
	playerUpdateListener listener.PlayerUpdateListener

	common Common
}

func NewPlayer(common Common) Player {
	return &player{
		player:               common.GetPlayer(),
		lastBattleLog:        common.GetLastBattleLog(),
		battleLog:            common.GetBattleLog(),
		common:               common,
		playerUpdateListener: common.GetPlayerUpdateListener(),
	}
}

func (b *player) GetDamageByWeapon(
	weapon string,
	player entities.Player,
	enemyDTO *battleDTO.EnemyDTO,
	gameCubeResultFirst uint,
	gameCubeResultSecond uint,
) (uint, error) {
	incrementDiceHits := uint(0)
	if enemyDTO.Abstract.MagicHit != nil {
		if !helpers.HasDebuff(enemyDTO.Instance.Debuff, entities.DebuffAliasMagicOffReason) {
			if b.common.Enemy().CountEnemyStep()%*enemyDTO.Abstract.MagicHit.Periodicity == 0 {
				if enemyDTO.Abstract.MagicHit.MinDiceHits != nil {
					result, err := expression.RunAndReturnRoundUint(
						fmt.Sprintf("%d %s", gameCubeResultFirst+gameCubeResultSecond, *enemyDTO.Abstract.MagicHit.MinDiceHits),
					)
					if err != nil {
						return 0, err
					}

					incrementDiceHits = result
				}
			}
		}
	}

	if weapon == battleDTO.Hand {
		if gameCubeResultFirst+gameCubeResultSecond > enemyDTO.Abstract.MinDiceHits+incrementDiceHits {
			return gameCubeResultFirst + gameCubeResultSecond - 6, nil
		}

		return 0, nil
	}

	for _, item := range player.Weapons {
		if weapon == item.Item {
			b.GetPlayerWeapon(weapon)

			if gameCubeResultFirst+gameCubeResultSecond > item.MinCubeHit {
				return item.Damage, nil
			}

			return 0, nil
		}
	}

	return 0, battleDTO.ErrBattleWeaponNotDefined
}

func (b *player) GetPlayerHealth(playerHealth uint, damage uint) uint {
	if playerHealth < damage {
		return 0
	}

	return playerHealth - damage
}

func (b *player) GetPlayerWeapon(battleWeapon string) []entities.Weapons {
	for index, playerWeapon := range b.player.Weapons {
		if playerWeapon.Item == battleWeapon {
			switch playerWeapon.Item {
			case battleDTO.Lightning:
				if playerWeapon.Count > 0 {
					b.player.Weapons[index].Count -= 1
				}
			case battleDTO.BallLightning:
				if playerWeapon.Count > 0 {
					b.player.Weapons[index].Count -= 1
				}
			}
		}
	}

	return b.player.Weapons
}

func (b *player) UpdatePlayerAfterStep() error {
	if b.player.Section.Type == battleDTO.SectionTypeSleepy {
		err := b.playerUpdateListener.Handle(context.Background(), event.PlayerUpdateEvent{
			PlayerID:  b.player.ID,
			SectionID: &b.player.SectionID,
			Health:    &b.player.Health,
			Weapons:   &b.player.Weapons,
			Meds:      &b.player.Meds,
			Bag:       &b.player.Bag,
			BonusList: &b.player.Bonus,
		})
		if err != nil {
			return err
		}
	}

	if b.player.Debuff != nil {
		for index, debuff := range b.player.Debuff {
			if debuff.Duration != nil {
				if *debuff.Duration-1 == 0 {
					b.player.Debuff = append(b.player.Debuff[:index], b.player.Debuff[index+1:]...)
				} else {
					*b.player.Debuff[index].Duration = *b.player.Debuff[index].Duration - 1
				}

			}
		}
	}

	if b.player.Buff != nil {
		for index, buff := range b.player.Buff {
			if buff.Duration != nil {
				if *buff.Duration-1 == 0 {
					b.player.Buff = append(b.player.Buff[:index], b.player.Buff[index+1:]...)
				} else {
					*b.player.Buff[index].Duration = *b.player.Buff[index].Duration - 1
				}

			}
		}
	}

	err := b.playerUpdateListener.Handle(context.Background(), event.PlayerUpdateEvent{
		PlayerID:  b.player.ID,
		SectionID: &b.player.SectionID,
		//ReturnToSection: &playerChanged.ReturnToSection,
		Health:    &b.player.Health,
		Weapons:   &b.player.Weapons,
		Meds:      &b.player.Meds,
		Bag:       &b.player.Bag,
		Debuff:    &b.player.Debuff,
		Buff:      &b.player.Buff,
		BonusList: &b.player.Bonus,
		//Description: fmt.Sprintf("%s", AntiPoisonSpellName),
	})
	if err != nil {
		return err
	}

	return nil
}
