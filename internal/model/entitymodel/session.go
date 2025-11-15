package entitymodel

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID            uuid.UUID  `db:"id"`
	Name          string     `db:"name"`
	DeckType      string     `db:"deck_type"`
	CardsRevealed bool       `db:"cards_revealed"`
	CreatorID     uuid.UUID  `db:"creator_id"`
	CreatorName   string     `db:"creator_name"`
	AllowEmoji    bool       `db:"allow_emoji"`
	AutoReveal    bool       `db:"auto_reveal"`
	CreatedVia    string     `db:"created_via"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}
