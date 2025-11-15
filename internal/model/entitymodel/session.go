package entitymodel

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID            uuid.UUID
	Name          string
	DeckType      string
	CardsRevealed bool
	CreatorID     uuid.UUID
	CreatorName   string
	AllowEmoji    bool
	AutoReveal    bool
	CreatedVia    string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
