package dto

import "context"

type (
	CustomSection interface {
		Handle(
			ctx context.Context,
		) (CustomSectionDTO, error)
	}

	CustomSectionDTO struct {
		SectionText string
		GotoItems   []TransitionDTO
	}
)
