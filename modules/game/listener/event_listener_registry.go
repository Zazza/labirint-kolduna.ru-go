package listener

import (
	"gamebook-backend/modules/game/dto"
	"gorm.io/gorm"
)

type EventListenerRegistry interface {
	Get(alias string) (EventListener, error)
}

type eventListenerRegistry struct {
	db        *gorm.DB
	factories map[string]EventListener
}

func NewEventListenerRegistry(db *gorm.DB) EventListenerRegistry {
	registry := &eventListenerRegistry{
		db:        db,
		factories: make(map[string]EventListener),
	}

	registry.RegisterDefaults()
	return registry
}

func (r *eventListenerRegistry) Get(alias string) (EventListener, error) {
	factory, exists := r.factories[alias]
	if !exists {
		return nil, dto.MessageEventListenerNotDefined
	}
	return factory, nil
}

func (r *eventListenerRegistry) RegisterDefaults() {
	r.factories["enemy_update"] = NewEnemyUpdateListener(r.db)
	r.factories["player_update"] = NewPlayerUpdateListener(r.db)
	r.factories["player_section"] = NewPlayerSectionListener(nil, r.db)
}
