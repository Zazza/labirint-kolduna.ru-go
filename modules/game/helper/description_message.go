package helper

import (
	"context"
	"gamebook-backend/modules/game/channel"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/pkg/helpers"

	"github.com/google/uuid"
)

var eventChannel channel.EventChannel

func InitEventChannel(ec channel.EventChannel) {
	eventChannel = ec
}

func DescriptionMessage(playerID uuid.UUID, message string) {
	escapedMessage := helpers.EscapeHTML(message)
	if eventChannel != nil {
		eventChannel.SendPlayerSectionUpdate(context.Background(), event.PlayerSectionEvent{
			PlayerID:    playerID,
			Description: &escapedMessage,
		})
	}
}

func DescriptionMessageWithContext(ctx context.Context, playerID uuid.UUID, message string) {
	escapedMessage := helpers.EscapeHTML(message)
	if eventChannel != nil {
		eventChannel.SendPlayerSectionUpdate(ctx, event.PlayerSectionEvent{
			PlayerID:    playerID,
			Description: &escapedMessage,
		})
	}
}

func SafeHTMLDescriptionMessage(playerID uuid.UUID, html string) {
	if eventChannel != nil {
		eventChannel.SendPlayerSectionUpdate(context.Background(), event.PlayerSectionEvent{
			PlayerID:    playerID,
			Description: &html,
		})
	}
}

func SafeHTMLDescriptionMessageWithContext(ctx context.Context, playerID uuid.UUID, html string) {
	if eventChannel != nil {
		eventChannel.SendPlayerSectionUpdate(ctx, event.PlayerSectionEvent{
			PlayerID:    playerID,
			Description: &html,
		})
	}
}
