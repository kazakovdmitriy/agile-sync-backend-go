package apimodel

import "time"

type Session struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	DeckType      string     `json:"deck_type"`
	CardsRevealed bool       `json:"cards_revealed"`
	CreatorID     string     `json:"creator_id"`
	CreatorName   string     `json:"creator_name"`
	AllowEmoji    bool       `json:"allow_emoji"`
	AutoReveal    bool       `json:"auto_reveal"`
	CreatedVia    string     `json:"created_via"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}
