package dto

import (
	"errors"
)

var (
	MessageBribeNotFoundSuccessTransition = errors.New("success transition not found")
	MessageBribeNotFoundFailTransition    = errors.New("fail transition not found")
	MessageBribeLogicError                = errors.New("bribe logic error")
)
