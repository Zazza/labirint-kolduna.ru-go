package section

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
)

type SectionManager interface {
	GetSection(ctx context.Context, sectionID string) (entities.Section, error)
	GetSectionByNumber(ctx context.Context, sectionNumber uint) (entities.Section, error)
	GetSectionsByNumbers(ctx context.Context, sectionNumbers []uint) ([]entities.Section, error)
	GetAvailableTransitions(ctx context.Context, section entities.Section, player entities.Player) ([]dto.TransitionDTO, error)
	ProcessTransition(ctx context.Context, player entities.Player, transitionID string) (dto.ActionResponse, error)
	GetSectionTemplate(ctx context.Context, section entities.Section, player entities.Player) (dto.SectionTemplate, error)
	CheckDiceRequirement(ctx context.Context, section entities.Section) bool
	IsCustomSection(ctx context.Context, number uint) bool
	GetSleepySection(ctx context.Context, number uint) (entities.Section, error)
}

type SectionValidator interface {
	ValidateTransition(ctx context.Context, section entities.Section, transition entities.Transition, player entities.Player) error
	ValidateSection(ctx context.Context, section entities.Section, player entities.Player) error
}

type SectionTemplateProvider interface {
	GetSectionTemplate(ctx context.Context, section entities.Section, player entities.Player) (dto.SectionTemplate, error)
}

type DiceRequirementChecker interface {
	IsDicesRequired(section entities.Section) bool
}

type TransitionProcessor interface {
	ProcessTransition(ctx context.Context, section entities.Section, transition entities.Transition, player entities.Player) (dto.ActionResponse, error)
}
