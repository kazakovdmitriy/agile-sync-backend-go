package converter

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
)

// SessionToUserSession конвертирует entity в API-модель
func SessionToUserSession(session *entitymodel.Session, userCount int) *apimodel.UserSessions {
	return &apimodel.UserSessions{
		ID:            session.ID,
		Name:          session.Name,
		DeckType:      session.DeckType,
		CardsRevealed: session.CardsRevealed,
		CreatorID:     session.CreatorID,
		CreatorName:   session.CreatorName,
		AllowEmoji:    session.AllowEmoji,
		AutoReveal:    session.AutoReveal,
		CreatedVia:    session.CreatedVia,
		CreatedAt:     session.CreatedAt,
		UpdatedAt:     session.UpdatedAt,
		UserCount:     userCount,
	}
}
