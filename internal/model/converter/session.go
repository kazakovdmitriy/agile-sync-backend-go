package converter

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/dbmodel"
	"backend_go/internal/model/entitymodel"
)

func SessionEntityToAPI(session *entitymodel.Session) *apimodel.Session {
	if session == nil {
		return nil
	}

	return &apimodel.Session{
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
	}
}

func SessionDBToEntity(session *dbmodel.Session) *entitymodel.Session {
	if session == nil {
		return nil
	}

	entitySession := &entitymodel.Session{
		ID:            session.ID,
		Name:          session.Name,
		DeckType:      session.DeckType,
		CardsRevealed: session.CardsRevealed,
		CreatorID:     session.CreatorID,
		CreatorName:   session.CreatorName,
		AllowEmoji:    session.AllowEmoji,
		AutoReveal:    session.AutoReveal,
		CreatedVia:    session.CreatedVia,
		CreatedAt:     &session.CreatedAt,
	}

	// Конвертируем UpdatedAt
	if session.UpdatedAt != nil {
		updatedAt := *session.UpdatedAt
		entitySession.UpdatedAt = &updatedAt
	}

	return entitySession
}

func SessionAPIToEntity(session *apimodel.Session) *entitymodel.Session {
	if session == nil {
		return nil
	}

	return &entitymodel.Session{
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
	}
}

func SessionEntityToDB(session *entitymodel.Session) *dbmodel.Session {
	if session == nil {
		return nil
	}

	dbSession := &dbmodel.Session{
		ID:            session.ID,
		Name:          session.Name,
		DeckType:      session.DeckType,
		CardsRevealed: session.CardsRevealed,
		CreatorID:     session.CreatorID,
		CreatorName:   session.CreatorName,
		AllowEmoji:    session.AllowEmoji,
		AutoReveal:    session.AutoReveal,
		CreatedVia:    session.CreatedVia,
	}

	// Конвертируем время
	if session.CreatedAt != nil {
		dbSession.CreatedAt = *session.CreatedAt
	}
	if session.UpdatedAt != nil {
		updatedAt := *session.UpdatedAt
		dbSession.UpdatedAt = &updatedAt
	}

	return dbSession
}
