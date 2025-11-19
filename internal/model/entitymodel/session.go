package entitymodel

import (
	"backend_go/internal/model"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID            uuid.UUID      `db:"id"            json:"id"`
	Name          string         `db:"name"          json:"name"`
	DeckType      model.DeckType `db:"deck_type"     json:"deck_type"`
	CardsRevealed bool           `db:"cards_revealed" json:"cards_revealed"`
	CreatorID     uuid.UUID      `db:"creator_id"    json:"creator_id"`
	CreatorName   string         `db:"creator_name"  json:"creator_name"`
	AllowEmoji    bool           `db:"allow_emoji"   json:"allow_emoji"`
	AutoReveal    bool           `db:"auto_reveal"   json:"auto_reveal"`
	CreatedVia    string         `db:"created_via"   json:"created_via"`
	CreatedAt     *time.Time     `db:"created_at"    json:"created_at"`
	UpdatedAt     *time.Time     `db:"updated_at"    json:"updated_at"`
}
