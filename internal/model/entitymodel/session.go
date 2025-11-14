package entitymodel

import "time"

type Session struct {
	ID            string
	Name          string
	DeckType      string
	CardsRevealed bool
	CreatorID     string
	CreatorName   string
	AllowEmoji    bool
	AutoReveal    bool
	CreatedVia    string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
