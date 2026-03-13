package channel

import _ "gamebook-backend/modules/game/listener/event"

type EventService interface {
	GetChannel() EventChannel
}

type eventService struct {
	channel EventChannel
}

func NewEventService(channel EventChannel) EventService {
	return &eventService{
		channel: channel,
	}
}

func (s *eventService) GetChannel() EventChannel {
	return s.channel
}
