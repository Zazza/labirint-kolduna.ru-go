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
	"math/rand"
	"strings"
)

type Enemy interface {
	GetEnemyHealth() uint
	GetAllEnemiesHealth() uint
	GetEnemy() (*battleDTO.EnemyDTO, error)
	GetEnemyDamage(enemy *battleDTO.EnemyDTO, gameCubeResultFirst uint, gameCubeResultSecond uint) (battleDTO.EnemyDamageDTO, error)
	CountEnemyStep() uint
	UpdateEnemyAfterStep(enemy *entities.PlayerSectionEnemy) error
}

type enemy struct {
	player              *entities.Player
	lastBattleLog       *entities.Battle
	battleLog           *[]entities.Battle
	dices               *entities.Dice
	enemies             *[]entities.PlayerSectionEnemy
	enemyUpdateListener listener.EnemyUpdateListener

	common Common
}

func NewEnemy(common Common) Enemy {
	return &enemy{
		player:              common.GetPlayer(),
		lastBattleLog:       common.GetLastBattleLog(),
		battleLog:           common.GetBattleLog(),
		enemies:             common.GetEnemies(),
		common:              common,
		enemyUpdateListener: common.GetEnemyUpdateListener(),
	}
}

func (b *enemy) GetEnemyHealth() uint {
	enemyDTO, err := b.GetEnemy()
	if err != nil {
		return 0
	}

	return enemyDTO.Instance.Health
}

func (b *enemy) GetAllEnemiesHealth() uint {
	enemyHealth := uint(0)
	for _, currentEnemy := range *b.enemies {
		enemyHealth += currentEnemy.Health
	}

	if enemyHealth < 0 {
		return 0
	}

	return enemyHealth
}

func (b *enemy) GetEnemy() (*battleDTO.EnemyDTO, error) {
	if b.player.Section.BattleSteps == nil {
		return nil, battleDTO.ErrStepNotDefined
	}

	nextStepIndex := uint(0)
	if b.lastBattleLog != nil {
		nextStepIndex = b.common.Step().GetNextStepIndex(b.lastBattleLog.Step)
	}

	var enemyAlias string

	steps, err := b.common.Step().GetCurrentBattleSteps()
	if err != nil {
		return nil, err
	}
	for {
		if steps[nextStepIndex] != AttackingPlayer {
			enemyAlias = steps[nextStepIndex]
			break
		}
		if nextStepIndex+1 == uint(len(steps)) {
			nextStepIndex = 0
		} else {
			nextStepIndex++
		}
	}

	var abstractEnemy entities.Enemy
	var instanceEnemy entities.PlayerSectionEnemy

	var abstractCompanion *entities.Enemy
	var instanceCompanion entities.PlayerSectionEnemy

	for _, enemySection := range b.player.Section.SectionEnemies {
		if enemySection == nil {
			continue
		}

		enemies := strings.Split(enemyAlias, "&&")
		var enemyCompanionAlias string
		if len(enemies) > 1 {
			enemyAlias = enemies[0]
			enemyCompanionAlias = enemies[1]

			if enemySection.Alias == enemyCompanionAlias {
				for _, instance := range *b.enemies {
					if instance.EnemyID == enemySection.ID {
						abstractCompanion = enemySection
						instanceCompanion = instance
					}
				}
			}
		}

		if enemySection.Alias == enemyAlias {
			for _, instance := range *b.enemies {
				if instance.EnemyID == enemySection.ID {
					abstractEnemy = *enemySection
					instanceEnemy = instance
				}
			}
		}
	}

	if instanceEnemy.Health == 0 {
		return &battleDTO.EnemyDTO{
			Abstract: abstractCompanion,
			Instance: &instanceCompanion,
		}, nil
	}

	if abstractCompanion != nil {
		return &battleDTO.EnemyDTO{
			Abstract: &abstractEnemy,
			Instance: &instanceEnemy,
			Companion: &battleDTO.EnemyDTO{
				Abstract: abstractCompanion,
				Instance: &instanceCompanion,
			},
		}, nil
	} else if &abstractEnemy != nil {
		return &battleDTO.EnemyDTO{
			Abstract: &abstractEnemy,
			Instance: &instanceEnemy,
		}, nil
	}

	return nil, battleDTO.ErrEnemyNotDefined
}

func (b *enemy) GetEnemyDamage(
	enemyDTO *battleDTO.EnemyDTO,
	gameCubeResultFirst uint,
	gameCubeResultSecond uint,
) (battleDTO.EnemyDamageDTO, error) {
	damage := uint(0)

	// TODO: Бафы и дебафы

	if helpers.HasDebuff(enemyDTO.Instance.Debuff, entities.DebuffAliasMagicOffReason) &&
		enemyDTO.Abstract.DamageType == entities.AliasMagicReason {
		return battleDTO.EnemyDamageDTO{
			Damage:      0,
			Description: "Действует антимагия",
		}, nil
	}

	// MagicHit
	if enemyDTO.Abstract.MagicHit != nil {
		if b.CountEnemyStep()%*enemyDTO.Abstract.MagicHit.Periodicity == 0 {
			if enemyDTO.Abstract.MagicHit.DicesValues != nil {
				result, err := expression.RunAndReturnBoolean(
					fmt.Sprintf("%d %s", gameCubeResultFirst+gameCubeResultSecond, *enemyDTO.Abstract.MagicHit.DicesValues),
				)
				if err != nil {
					return battleDTO.EnemyDamageDTO{}, err
				}

				if result && *enemyDTO.Abstract.MagicHit.InstantKill {
					b.player.Health = 0
				}

				if enemyDTO.Abstract.MagicHit.Damage != nil {
					damage = *enemyDTO.Abstract.MagicHit.Damage
				}

				// TODO: другие, например -health?
			}
		}
	} else if enemyDTO.Instance.Weapons != nil {
		var weapon *entities.EnemyWeapon
		index := rand.Intn(len(enemyDTO.Instance.Weapons))
		weapon = &enemyDTO.Instance.Weapons[index]

		damage = weapon.Damage

		if weapon.Count != nil {
			*weapon.Count--
		}
	} else {
		if gameCubeResultFirst+gameCubeResultSecond+enemyDTO.Abstract.Damage > enemyDTO.Abstract.MinDiceHits {
			damage = gameCubeResultFirst + gameCubeResultSecond - enemyDTO.Abstract.MinDiceHits + enemyDTO.Abstract.Damage
		}
	}

	if damage == 0 {
		return battleDTO.EnemyDamageDTO{
			Damage:      0,
			Description: "Враг промахнулся",
		}, nil
	}

	if enemyDTO.Abstract.PlayerArmor {
		if damage <= ChainMailProtection {
			return battleDTO.EnemyDamageDTO{
				Damage:      0,
				Description: "Броня сдержала урон",
			}, nil
		}

		damage -= ChainMailProtection
	}

	playerHealth := b.player.Health - damage
	if b.player.Health < damage {
		playerHealth = 0
	}

	return battleDTO.EnemyDamageDTO{
		Damage: damage,
		Description: fmt.Sprintf(
			"%s нанес %d урона. Теперь у тебя %d здоровья",
			enemyDTO.Abstract.Name,
			damage,
			playerHealth,
		),
	}, nil
}

func (b *enemy) CountEnemyStep() uint {
	count := uint(0)
	for _, item := range *b.battleLog {
		if item.Attacking == AttackingEnemy {
			count++
		}
	}

	return count
}

func (b *enemy) UpdateEnemyAfterStep(enemy *entities.PlayerSectionEnemy) error {
	for elIndex, debuff := range enemy.Debuff {
		if debuff.Duration != nil {
			if *debuff.Duration-1 == 0 {
				enemy.Debuff = append(enemy.Debuff[:elIndex], enemy.Debuff[elIndex+1:]...)
			} else {
				*enemy.Debuff[elIndex].Duration--
			}
		}
	}

	for elIndex, buff := range enemy.Buff {
		if buff.Duration != nil {
			if *buff.Duration-1 == 0 {
				enemy.Buff = append(enemy.Buff[:elIndex], enemy.Buff[elIndex+1:]...)
			} else {
				*enemy.Buff[elIndex].Duration--
			}
		}
	}

	err := b.enemyUpdateListener.Handle(context.Background(), event.EnemyUpdateEvent{
		PlayerID:  b.player.ID,
		SectionID: b.player.SectionID,
		EnemyID:   enemy.EnemyID,
		Buff:      &enemy.Buff,
		Debuff:    &enemy.Debuff,
		Health:    &enemy.Health,
		Weapons:   &enemy.Weapons,
	})
	if err != nil {
		return err
	}

	return nil
}
