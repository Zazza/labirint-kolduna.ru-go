package bonus

import (
	"context"
	"gamebook-backend/modules/game/dto"
)

type BonusOptionHandler func(ctx context.Context, req dto.BonusRequest) error

type OptionHandlerMap map[string]BonusOptionHandler

func NewOptionHandlerMap(handlers OptionHandlerMap) OptionHandlerMap {
	return handlers
}

func (m OptionHandlerMap) Execute(ctx context.Context, req dto.BonusRequest) error {
	handler, exists := m[req.Option]
	if !exists {
		return dto.ErrBattleNotFound
	}
	return handler(ctx, req)
}

type MagicRingHandlers struct {
	leftOption  BonusOptionHandler
	rightOption BonusOptionHandler
}

func NewMagicRingHandlers(left, right BonusOptionHandler) *MagicRingHandlers {
	return &MagicRingHandlers{
		leftOption:  left,
		rightOption: right,
	}
}

func (h *MagicRingHandlers) ToMap() OptionHandlerMap {
	return map[string]BonusOptionHandler{
		MagicRingOptions[0]: h.leftOption,
		MagicRingOptions[1]: h.rightOption,
	}
}

type MagicDuckHandlers struct {
	antiMagicOption BonusOptionHandler
	sectionOption   BonusOptionHandler
}

func NewMagicDuckHandlers(antiMagic, section BonusOptionHandler) *MagicDuckHandlers {
	return &MagicDuckHandlers{
		antiMagicOption: antiMagic,
		sectionOption:   section,
	}
}

func (h *MagicDuckHandlers) ToMap() OptionHandlerMap {
	return map[string]BonusOptionHandler{
		MagicDuckOptions[0]: h.antiMagicOption,
		MagicDuckOptions[1]: h.sectionOption,
	}
}
