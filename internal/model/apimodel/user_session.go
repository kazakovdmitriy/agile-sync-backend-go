package apimodel

import (
	"backend_go/internal/model"
	"github.com/google/uuid"
	"time"
)

// Модель для списка пользовательских сессий
type UserSessions struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	DeckType      model.DeckType `json:"deck_type"`
	CardsRevealed bool           `json:"cards_revealed"`
	CreatorID     uuid.UUID      `json:"creator_id"`
	CreatorName   string         `json:"creator_name"`
	AllowEmoji    bool           `json:"allow_emoji"`
	AutoReveal    bool           `json:"auto_reveal"`
	CreatedVia    string         `json:"created_via"`
	CreatedAt     *time.Time     `json:"created_at,omitempty"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
	UserCount     int            `json:"user_count"`
}
