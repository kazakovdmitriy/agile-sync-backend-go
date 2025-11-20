package apimodel

import (
	"backend_go/internal/model"
	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID            `json:"id"`
	Name           string               `json:"name"`
	DeckType       model.DeckType       `json:"deck_type"`
	CardsRevealed  bool                 `json:"cards_revealed"`
	CreatorID      uuid.UUID            `json:"creator_id"`
	CreatorName    string               `json:"creator_name"`
	AutoReveal     bool                 `json:"auto_reveal"`
	AllowEmoji     bool                 `json:"allow_emoji"`
	CreatedVia     string               `json:"created_via"`
	DeckValues     []string             `json:"deck_values"`
	Users          []UsersInSession     `json:"users"`
	Votes          map[uuid.UUID]string `json:"votes"`
	AverageVote    *string              `json:"average_vote"`
	MostCommonVote *string              `json:"most_common_vote"`
}

type UsersInSession struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsCreator bool      `json:"is_creator"`
	IsWatcher bool      `json:"is_watcher"`
}
