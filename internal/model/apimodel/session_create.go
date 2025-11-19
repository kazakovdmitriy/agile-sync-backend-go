package apimodel

import "backend_go/internal/model"

type SessionCreate struct {
	Name     string         `json:"name"`
	DeckType model.DeckType `json:"deck_type"`
}
