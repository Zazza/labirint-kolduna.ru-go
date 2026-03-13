package listener

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type EnemyUpdateListener interface {
	Handle(ctx context.Context, e event.Event) error
}

type enemyListener struct {
	db                           *gorm.DB
	playerSectionEnemyRepository repository.PlayerSectionEnemyRepository
}

func NewEnemyUpdateListener(
	db *gorm.DB,
) EnemyUpdateListener {
	playerSectionEnemyRepository := repository.NewPlayerSectionEnemyRepository(db)

	return &enemyListener{
		db:                           db,
		playerSectionEnemyRepository: playerSectionEnemyRepository,
	}
}

func (l *enemyListener) Handle(ctx context.Context, e event.Event) error {
	eventEnemyUpdate, ok := e.(event.EnemyUpdateEvent)
	if !ok {
		return nil
	}

	enemies, err := l.playerSectionEnemyRepository.GetEnemiesByPlayerIDAndSectionID(
		ctx,
		l.db,
		eventEnemyUpdate.PlayerID,
		eventEnemyUpdate.SectionID,
	)
	if err != nil {
		return err
	}

	enemyChanged, err := l.enemyUpdateActions(ctx, *enemies, eventEnemyUpdate)
	if err != nil {
		return err
	}

	return l.playerSectionEnemyRepository.UpdateAll(ctx, l.db, enemyChanged)
}

func (l *enemyListener) enemyUpdateActions(
	ctx context.Context,
	enemies []entities.PlayerSectionEnemy,
	eventEnemyUpdate event.EnemyUpdateEvent,
) ([]entities.PlayerSectionEnemy, error) {
	if eventEnemyUpdate.Health != nil {
		for index, item := range enemies {
			if item.EnemyID == eventEnemyUpdate.EnemyID {
				enemies[index].Health = *eventEnemyUpdate.Health
			}
		}
	}

	if eventEnemyUpdate.Debuff != nil {
		for index, item := range enemies {
			if item.EnemyID == eventEnemyUpdate.EnemyID {
				enemies[index].Debuff = *eventEnemyUpdate.Debuff
			}
		}
	}
	if eventEnemyUpdate.Buff != nil {
		for index, item := range enemies {
			if item.EnemyID == eventEnemyUpdate.EnemyID {
				enemies[index].Buff = *eventEnemyUpdate.Buff
			}
		}
	}
	if eventEnemyUpdate.Weapons != nil {
		for index, item := range enemies {
			if item.EnemyID == eventEnemyUpdate.EnemyID {
				enemies[index].Weapons = *eventEnemyUpdate.Weapons
			}
		}
	}

	return enemies, nil
}
