package ability

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"
	"github.com/google/uuid"
)

type BonusLogic struct {
	diceRepository       repository.DiceRepository
	playerRepository     repository.PlayerRepository
	playerUpdateListener listener.PlayerUpdateListener
	bonusRepository      repository.BonusRepository
}

func NewBonusLogic(
	diceRepo repository.DiceRepository,
	playerRepo repository.PlayerRepository,
	playerUpdateListener listener.PlayerUpdateListener,
	bonusRepo repository.BonusRepository,
) *BonusLogic {
	return &BonusLogic{
		diceRepository:       diceRepo,
		playerRepository:     playerRepo,
		playerUpdateListener: playerUpdateListener,
		bonusRepository:      bonusRepo,
	}
}

func (b *BonusLogic) Validate(ctx context.Context, player entities.Player) error {
	if player.Section.Type == dto.SectionTypeSleepy {
		return fmt.Errorf("cannot use bonus in sleep section type %s", player.Section.Type)
	}

	if len(player.Bonus) == 0 {
		return fmt.Errorf("player has no bonuses")
	}

	return nil
}

func (b *BonusLogic) Execute(ctx context.Context, player entities.Player, req dto.BonusRequest) error {
	if player.Section.Type == dto.SectionTypeSleepy {
		return fmt.Errorf("cannot use bonus in sleep section type %s", player.Section.Type)
	}

	if len(player.Bonus) == 0 && req.Bonus == "" {
		return fmt.Errorf("player has no bonuses to use, cannot use empty alias")
	}

	usedBonus, err := b.bonusRepository.GetByPlayerID(ctx, player.ID)
	if err != nil {
		return fmt.Errorf("bonus not found for player: %w", err)
	}

	_ = b.playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID: player.ID,
		Bonus:    &usedBonus,
	})

	b.sendDescription(ctx, player.ID, fmt.Sprintf("Bonus used: %s", req.Bonus))
	return nil
}

func (b *BonusLogic) Result() dto.BonusDTO {
	return dto.BonusDTO{
		Success: true,
		Message: "Bonus used successfully",
	}
}

func (b *BonusLogic) sendDescription(ctx context.Context, playerID uuid.UUID, description string) {
	helper.SafeHTMLDescriptionMessageWithContext(ctx, playerID, fmt.Sprintf("<p>✨ Использован бонус: %s</p>", description))
}
