package section

import (
	"context"
	"errors"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/expression"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener/event"
	player2 "gamebook-backend/modules/game/player"
	"gamebook-backend/modules/game/repository"
	"log"

	"gorm.io/gorm"
)

type SectionUpdate interface {
	Update(ctx context.Context, player entities.Player, transition entities.Transition) error
}

type sectionUpdate struct {
	db                           *gorm.DB
	playerRepository             repository.PlayerRepository
	sectionRepository            repository.SectionRepository
	diceRepository               repository.DiceRepository
	playerSectionRepository      repository.PlayerSectionRepository
	playerSectionEnemyRepository repository.PlayerSectionEnemyRepository
}

func NewSectionUpdate(
	db *gorm.DB,
) SectionUpdate {
	playerRepo := repository.NewPlayerRepository(db)
	sectionRepo := repository.NewSectionRepository(db)
	diceRepo := repository.NewDiceRepository(db)
	playerSectionRepository := repository.NewPlayerSectionRepository(db)
	playerSectionEnemyRepo := repository.NewPlayerSectionEnemyRepository(db)
	return &sectionUpdate{
		db:                           db,
		playerRepository:             playerRepo,
		sectionRepository:            sectionRepo,
		diceRepository:               diceRepo,
		playerSectionRepository:      playerSectionRepository,
		playerSectionEnemyRepository: playerSectionEnemyRepo,
	}
}

func NewSectionUpdateWithRepositories(
	playerRepo repository.PlayerRepository,
	sectionRepo repository.SectionRepository,
	diceRepo repository.DiceRepository,
	playerSectionRepo repository.PlayerSectionRepository,
	playerSectionEnemyRepo repository.PlayerSectionEnemyRepository,
) SectionUpdate {
	return &sectionUpdate{
		playerRepository:             playerRepo,
		sectionRepository:            sectionRepo,
		diceRepository:               diceRepo,
		playerSectionRepository:      playerSectionRepo,
		playerSectionEnemyRepository: playerSectionEnemyRepo,
	}
}

func (l *sectionUpdate) Update(ctx context.Context, player entities.Player, transition entities.Transition) error {
	enemies, err := l.playerSectionEnemyRepository.GetEnemiesByPlayerIDAndSectionID(
		ctx,
		l.db,
		player.ID,
		player.SectionID,
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	changedEnemies := l.sectionEndChangeEnemiesActions(ctx, enemies)
	for _, enemy := range *changedEnemies {
		err = l.playerSectionEnemyRepository.Update(ctx, l.db, enemy)
		if err != nil {
			return err
		}
	}

	playerChanged := l.sectionEndChangePlayerActions(ctx, &player)

	err = l.diceRepository.Remove(
		ctx,
		l.db,
		player.ID,
	)
	if err != nil {
		return err
	}

	playerChanged, err = l.sectionStartActions(ctx, playerChanged, transition)
	if err != nil {
		log.Println("[sectionListener] update section err:", err)
	}

	playerUpdate := player2.NewPlayerUpdate(l.db, player.ID)
	_, err = playerUpdate.Update(ctx, event.PlayerUpdateEvent{
		PlayerID:  playerChanged.ID,
		SectionID: &playerChanged.SectionID,
		Health:    &playerChanged.Health,
		HealthMax: &playerChanged.HealthMax,
		Weapons:   &playerChanged.Weapons,
		Meds:      &playerChanged.Meds,
		Bag:       &playerChanged.Bag,
		Debuff:    &playerChanged.Debuff,
		Buff:      &playerChanged.Buff,
		BonusList: &playerChanged.Bonus,
		Gold:      &playerChanged.Gold,
	})
	if err != nil {
		return err
	}

	return nil
}

func (l *sectionUpdate) sectionEndChangeEnemiesActions(
	ctx context.Context,
	enemies *[]entities.PlayerSectionEnemy,
) *[]entities.PlayerSectionEnemy {
	var result []entities.PlayerSectionEnemy
	if enemies != nil {
		for _, enemy := range *enemies {
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

			result = append(result, enemy)
		}
	}

	return &result
}

func (l *sectionUpdate) sectionEndChangePlayerActions(
	ctx context.Context,
	player *entities.Player,
) *entities.Player {
	if player.Debuff != nil {
		for index, debuff := range player.Debuff {
			if debuff.Duration != nil {
				if *debuff.Duration-1 == 0 {
					player.Debuff = append(player.Debuff[:index], player.Debuff[index+1:]...)
				} else {
					*player.Debuff[index].Duration = *player.Debuff[index].Duration - 1
				}

			}
		}
	}

	if player.Buff != nil {
		for index, buff := range player.Buff {
			if buff.Duration != nil {
				if *buff.Duration-1 == 0 {
					player.Buff = append(player.Buff[:index], player.Buff[index+1:]...)
				} else {
					*player.Buff[index].Duration = *player.Buff[index].Duration - 1
				}

			}
		}
	}

	return player
}

func (l *sectionUpdate) sectionStartActions(
	ctx context.Context,
	player *entities.Player,
	transition entities.Transition,
) (*entities.Player, error) {
	changeDTO := ChangeDTO{
		Player:  player,
		Message: make([]string, 0),
	}

	var err error

	if transition.PlayerChange != nil {
		changeDTO, err = Change(ctx, transition, *player, l.db)
		if err != nil {
			return nil, err
		}

		player = changeDTO.Player
	}
	if transition.PlayerDebuff != nil {
		for _, debuff := range transition.PlayerDebuff {
			if debuff.Health != nil {
				if player.Health <= *debuff.Health {
					player.Debuff = append(player.Debuff, entities.Debuff{
						Alias:  *debuff.Alias,
						Health: debuff.Health,
					})
				} else {
					player.Health -= *debuff.Health

					player.Debuff = append(player.Debuff, entities.Debuff{
						Alias:  *debuff.Alias,
						Health: debuff.Health,
					})
				}
			}
			if debuff.HealthString != nil {
				value, err := expression.RunAndReturnRoundUint(fmt.Sprintf(
					"%d %s",
					player.Health,
					*debuff.HealthString,
				))
				if err != nil {
					return nil, err
				}

				decrementHealth := player.Health - value

				player.Health = player.Health - decrementHealth

				player.Debuff = append(player.Debuff, entities.Debuff{
					Alias:  *debuff.Alias,
					Health: &decrementHealth,
				})
			}
			if debuff.MinCubeHit != nil {
				player.Debuff = append(player.Debuff, entities.Debuff{
					MinCubeHit: debuff.MinCubeHit,
					Duration:   debuff.Duration,
				})
			}
		}
	}

	if transition.PlayerBuff != nil {
		for _, buff := range transition.PlayerBuff {
			if buff.Health != nil {
				player.Health += *buff.Health
			}
		}
	}

	player.SectionID = transition.TargetSectionID

	err = l.CreatePlayerSection(ctx, *changeDTO.Player)
	if err != nil {
		return nil, err
	}

	for _, msg := range changeDTO.Message {
		helper.DescriptionMessage(player.ID, msg)
	}

	return player, nil
}

func (l *sectionUpdate) CreatePlayerSection(ctx context.Context, player entities.Player) error {
	playerSection, err := l.playerSectionRepository.GetLastPlayerSection(ctx, l.db, player.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := l.playerSectionRepository.Create(
			ctx,
			l.db,
			player.ID,
			player.SectionID,
		)
		if err != nil {
			return err
		}

		return nil
	} else if err != nil {
		return err
	}
	if playerSection.SectionID != player.SectionID {
		err := l.playerSectionRepository.Create(
			ctx,
			l.db,
			player.ID,
			player.SectionID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
