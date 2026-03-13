package battle

import (
	"context"
	"errors"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	battleDTO "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/log"
	"gamebook-backend/modules/game/repository"
	template2 "gamebook-backend/modules/game/template"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Common interface {
	Action(weapon *string) (entities.Battle, error)

	Step() Step
	Enemy() Enemy
	Player() Player
	GetPlayerUpdateListener() listener.PlayerUpdateListener
	GetEnemyUpdateListener() listener.EnemyUpdateListener

	GetGotoSectionWin() ([]entities.Transition, error)
	GetGotoSectionLose() (uuid.UUID, error)

	GetPlayer() *entities.Player
	GetBattleLog() *[]entities.Battle
	GetLastBattleLog() *entities.Battle
	GetEnemies() *[]entities.PlayerSectionEnemy
	IsWin() bool
}

type battleCommon struct {
	player               *entities.Player
	lastBattleLog        *entities.Battle
	battleLog            *[]entities.Battle
	rollTheDicesPlayer   dice.RollTheDices
	rollTheDicesEnemy    dice.RollTheDices
	enemies              *[]entities.PlayerSectionEnemy
	battlesRepository    repository.BattleRepository
	logService           log.PlayerLogService
	enemyUpdateListener  listener.EnemyUpdateListener
	playerUpdateListener listener.PlayerUpdateListener
}

func NewCommon(
	ctx context.Context,
	db *gorm.DB,
	player *entities.Player,
) (Common, error) {
	battlesRepository := repository.NewBattleRepository(db)
	playerSectionEnemyRepository := repository.NewPlayerSectionEnemyRepository(db)

	return NewCommonWithRepositories(ctx, db, player, battlesRepository, playerSectionEnemyRepository, nil)
}

func NewCommonWithRepositories(
	ctx context.Context,
	db *gorm.DB,
	player *entities.Player,
	battlesRepository repository.BattleRepository,
	playerSectionEnemyRepository repository.PlayerSectionEnemyRepository,
	logService log.PlayerLogService,
) (Common, error) {
	battleLog, err := battlesRepository.GetByPlayerIdAndSectionNumber(db, player.ID, player.Section.Number)
	if err != nil {
		return nil, err
	}

	battleLogLast, err := battlesRepository.FindLastByPlayerIdAndSectionNumber(db, player.ID, player.Section.Number)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	enemies, err := playerSectionEnemyRepository.GetEnemiesByPlayerIDAndSectionID(ctx, db, player.ID, player.SectionID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		newEnemies, err := playerSectionEnemyRepository.Create(ctx, db, player.ID, player.SectionID, player.Section.SectionEnemies)
		if err != nil {
			return nil, err
		}

		enemies = &newEnemies
	}

	rollTheDicesPlayer := dice.NewRollTheDices(db, player)
	rollTheDicesEnemy := dice.NewRollTheDices(db, nil)

	enemyUpdateListener, err := listener.HandleEvent(db, "enemy_update")
	if err != nil {
		return nil, err
	}

	playerUpdateListener, err := listener.HandleEvent(db, "player_update")
	if err != nil {
		return nil, err
	}

	return &battleCommon{
		player:               player,
		lastBattleLog:        battleLogLast,
		battleLog:            &battleLog,
		rollTheDicesPlayer:   rollTheDicesPlayer,
		rollTheDicesEnemy:    rollTheDicesEnemy,
		enemies:              enemies,
		battlesRepository:    battlesRepository,
		logService:           logService,
		enemyUpdateListener:  enemyUpdateListener,
		playerUpdateListener: playerUpdateListener,
	}, nil
}

func (b *battleCommon) GetPlayer() *entities.Player {
	return b.player
}

func (b *battleCommon) GetLastBattleLog() *entities.Battle {
	return b.lastBattleLog
}

func (b *battleCommon) GetBattleLog() *[]entities.Battle {
	return b.battleLog
}

func (b *battleCommon) GetEnemies() *[]entities.PlayerSectionEnemy {
	return b.enemies
}

func (b *battleCommon) Step() Step {
	return NewStep(b)
}

func (b *battleCommon) Enemy() Enemy {
	return NewEnemy(b)
}

func (b *battleCommon) Player() Player {
	return NewPlayer(b)
}

func (b *battleCommon) GetPlayerUpdateListener() listener.PlayerUpdateListener {
	return b.playerUpdateListener
}

func (b *battleCommon) GetEnemyUpdateListener() listener.EnemyUpdateListener {
	return b.enemyUpdateListener
}

func (b *battleCommon) Action(weapon *string) (entities.Battle, error) {
	isMyMove, err := b.Step().IsMyMove()
	if err != nil {
		return entities.Battle{}, err
	}

	var step uint
	if b.lastBattleLog == nil {
		for index, item := range b.player.Section.BattleSteps {
			if *item == AttackingPlayer && isMyMove {
				step = uint(index)
				break
			}
			if *item == AttackingEnemy && !isMyMove {
				step = uint(index)
				break
			}
		}
	} else {
		step = b.lastBattleLog.Step + 1
	}

	steps, err := b.Step().GetCurrentBattleSteps()
	if err != nil {
		return entities.Battle{}, err
	}

	if step >= uint(len(steps)) {
		step = 0
	}

	var battle entities.Battle
	var damage uint
	var damageDTO battleDTO.WeaponDamageDto

	enemyDTO, err := b.Enemy().GetEnemy()
	if err != nil {
		return entities.Battle{}, err
	}

	if isMyMove {
		gameCubeResultFirst, gameCubeResultSecond, err := b.rollTheDicesPlayer.RollTheDices(context.Background(), *b.player)
		if err != nil {
			return entities.Battle{}, err
		}

		if enemyDTO.Abstract.OnlyDiceHits != nil {
			hit := 0
			description := "Ты промахнулся"
			if *gameCubeResultFirst+*gameCubeResultSecond >= *enemyDTO.Abstract.OnlyDiceHits {
				hit = 1
				description = fmt.Sprintf("Ты убил %s", enemyDTO.Abstract.Name)
			}

			damageDTO = battleDTO.WeaponDamageDto{
				Damage:      uint(hit),
				Description: description,
			}
		} else {
			damage, err = b.Player().GetDamageByWeapon(
				*weapon,
				*b.player,
				enemyDTO,
				*gameCubeResultFirst,
				*gameCubeResultSecond,
			)
			if err != nil {
				return entities.Battle{}, err
			}

			enemyHealth := b.Enemy().GetEnemyHealth()
			if damage > enemyHealth {
				enemyHealth = 0
			} else {
				enemyHealth -= damage
			}

			damageDTO = battleDTO.WeaponDamageDto{
				Damage: damage,
				Description: fmt.Sprintf(
					"Бил ты (%s) и нанес %d урона. Теперь у %s %d здоровья",
					Weapon[*weapon],
					damage,
					enemyDTO.Abstract.Name,
					enemyHealth,
				),
			}
		}

		enemyHealth := enemyDTO.Instance.Health - damageDTO.Damage
		if enemyDTO.Instance.Health < damageDTO.Damage {
			enemyHealth = 0
		}

		err = b.enemyUpdateListener.Handle(context.Background(), event.EnemyUpdateEvent{
			PlayerID:  b.player.ID,
			SectionID: b.player.SectionID,
			EnemyID:   enemyDTO.Instance.EnemyID,
			Health:    &enemyHealth,
		})
		if err != nil {
			return entities.Battle{}, err
		}

		if b.logService != nil {
			b.logService.LogBattleHit(b.player.ID, "player", enemyDTO.Abstract.Name, damageDTO.Damage, *weapon, nil, nil)
		}

		err = b.Player().UpdatePlayerAfterStep()
		if err != nil {
			return entities.Battle{}, err
		}

		battle = entities.Battle{
			Section:     b.player.Section.Number,
			PlayerID:    b.player.ID,
			Attacking:   AttackingPlayer,
			Damage:      damageDTO.Damage,
			Dice1:       *gameCubeResultFirst,
			Dice2:       *gameCubeResultSecond,
			Description: damageDTO.Description,
			Weapon:      *weapon,
			Step:        step,
		}

	} else {
		gameCubeResultFirst, gameCubeResultSecond, err := b.rollTheDicesEnemy.RollTheDices(context.Background(), *b.player)
		if err != nil {
			return entities.Battle{}, err
		}

		if enemyDTO.Instance.Debuff != nil {
			for _, item := range enemyDTO.Instance.Debuff {
				if item.Alias == entities.DebuffAliasSkipReason {
					damage := 0
					description := fmt.Sprintf("На %s действует заморозка", enemyDTO.Abstract.Name)

					err = b.Enemy().UpdateEnemyAfterStep(enemyDTO.Instance)
					if err != nil {
						return entities.Battle{}, err
					}

					battle = entities.Battle{
						Section:     b.player.Section.Number,
						PlayerID:    b.player.ID,
						Attacking:   AttackingEnemy,
						Damage:      uint(damage),
						Dice1:       *gameCubeResultFirst,
						Dice2:       *gameCubeResultSecond,
						Description: description,
						Step:        step,
					}

					return battle, nil
				}
			}
		}

		var enemyDamageCompanion battleDTO.EnemyDamageDTO
		var dicesCompanionTemplate string
		if enemyDTO.Companion != nil {
			gameCubeCompanionResultFirst, gameCubeCompanionResultSecond, err := b.rollTheDicesEnemy.RollTheDices(context.Background(), *b.player)
			if err != nil {
				return entities.Battle{}, err
			}

			dicesCompanionTemplate, err = template2.GetDicesTemplate(
				context.Background(),
				*gameCubeCompanionResultFirst,
				*gameCubeCompanionResultSecond,
				false,
			)
			if err != nil {
				return entities.Battle{}, err
			}

			enemyDamageCompanion, err = b.Enemy().GetEnemyDamage(enemyDTO.Companion, *gameCubeCompanionResultFirst, *gameCubeCompanionResultSecond)
			if err != nil {
				return entities.Battle{}, err
			}
		}

		enemyDamage, err := b.Enemy().GetEnemyDamage(enemyDTO, *gameCubeResultFirst, *gameCubeResultSecond)
		if err != nil {
			return entities.Battle{}, err
		}

		if &enemyDamageCompanion != nil {
			enemyDamage = battleDTO.EnemyDamageDTO{
				Damage:      enemyDamage.Damage + enemyDamageCompanion.Damage,
				Description: enemyDamage.Description + "<div class='section-log'>" + dicesCompanionTemplate + enemyDamageCompanion.Description + "</div>",
			}
		} else {
			enemyDamage = battleDTO.EnemyDamageDTO{
				Damage:      enemyDamage.Damage,
				Description: enemyDamage.Description,
			}
		}

		if (enemyDTO.Abstract.DamageType == entities.AliasPoisonReason) || (enemyDTO.Abstract.DamageType == entities.AliasMagicReason) {
			b.player.Debuff = append(b.player.Debuff, entities.Debuff{
				Alias:  enemyDTO.Abstract.DamageType,
				Health: &enemyDamage.Damage,
			})
		}

		b.player.Health = b.Player().GetPlayerHealth(b.player.Health, enemyDamage.Damage)

		err = b.playerUpdateListener.Handle(context.Background(), event.PlayerUpdateEvent{
			PlayerID: b.player.ID,
			Health:   &b.player.Health,
		})
		if err != nil {
			return entities.Battle{}, err
		}

		if b.logService != nil {
			b.logService.LogBattleHit(b.player.ID, enemyDTO.Abstract.Name, "player", enemyDamage.Damage, "", nil, nil)
		}

		err = b.Enemy().UpdateEnemyAfterStep(enemyDTO.Instance)
		if err != nil {
			return entities.Battle{}, err
		}

		battle = entities.Battle{
			Section:     b.player.Section.Number,
			PlayerID:    b.player.ID,
			Attacking:   AttackingEnemy,
			Damage:      damage,
			Dice1:       *gameCubeResultFirst,
			Dice2:       *gameCubeResultSecond,
			Description: enemyDamage.Description,
			Step:        step,
		}
	}

	return battle, nil
}

func (b *battleCommon) GetGotoSectionWin() ([]entities.Transition, error) {
	var result []entities.Transition

	if b.player.Section.Type == battleDTO.SectionTypeSleepy {
		result = append(result, entities.Transition{
			ID:   uuid.New(),
			Text: "Далее",
		})
	} else {
		for _, item := range b.player.Section.Transitions {
			if item.IsBattleWin != nil {
				if *item.IsBattleWin {
					result = append(result, item)
				}
			}
		}
	}
	if len(result) > 0 {
		return result, nil
	}

	return result, battleDTO.ErrSectionNotFound
}

func (b *battleCommon) GetGotoSectionLose() (uuid.UUID, error) {
	if b.player.Section.Type == battleDTO.SectionTypeSleepy {
		return uuid.New(), nil
	}

	for _, item := range b.player.Section.Transitions {
		if !*item.IsBattleWin {
			return item.ID, nil
		}
	}

	return uuid.UUID{}, battleDTO.ErrSectionNotFound
}

func (b *battleCommon) IsWin() bool {
	for _, item := range *b.enemies {
		if 0 < item.Health {
			return false
		}
	}

	return true
}
