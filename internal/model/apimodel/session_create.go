package apimodel

type SessionCreate struct {
	Name     string `json:"name"`
	DeckType string `json:"deck_type"`
}
