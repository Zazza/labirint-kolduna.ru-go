package dto

import "github.com/google/uuid"

type SectionNode struct {
	ID        uuid.UUID `json:"id"`
	Number    uint      `json:"number"`
	Title     string    `json:"title"`
	IsVisited bool      `json:"is_visited"`
	IsCurrent bool      `json:"is_current"`
}

type TransitionEdge struct {
	FromSection uint   `json:"from_section"`
	ToSection   uint   `json:"to_section"`
	Text        string `json:"text"`
	IsAvailable bool   `json:"is_available"`
}

type MapResponse struct {
	Sections    []SectionNode    `json:"sections"`
	Transitions []TransitionEdge `json:"transitions"`
}
