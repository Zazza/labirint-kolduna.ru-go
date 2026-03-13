package listener

import (
	"context"
	"gamebook-backend/modules/game/listener/event"

	"gorm.io/gorm"
)

type EventListenerFactory interface {
	GetEventListener(eventName string) (EventListener, error)
}

type eventListenerFactory struct {
	registry EventListenerRegistry
}

type EventListener interface {
	Handle(ctx context.Context, e event.Event) error
}

func NewEventListenerFactory(
	db *gorm.DB,
) EventListenerFactory {
	return &eventListenerFactory{
		registry: NewEventListenerRegistry(db),
	}
}

func (f *eventListenerFactory) GetEventListener(eventName string) (EventListener, error) {
	eventListener, err := f.registry.Get(eventName)
	if err != nil {
		return nil, err
	}

	return eventListener, err
}

func HandleEvent(db *gorm.DB, eventName string) (EventListener, error) {
	eventListener := NewEventListenerFactory(db)
	listener, err := eventListener.GetEventListener(eventName)
	if err != nil {
		return nil, err
	}
	return listener, nil
}
