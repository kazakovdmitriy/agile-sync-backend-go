package converter

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/dbmodel"
	"backend_go/internal/model/entitymodel"
)

func ReactionEntityToAPI(reaction *entitymodel.Reaction) *apimodel.Reaction {
	if reaction == nil {
		return nil
	}

	return &apimodel.Reaction{
		ID:         reaction.ID,
		SessionID:  reaction.SessionID,
		FromUserID: reaction.FromUserID,
		ToUserID:   reaction.ToUserID,
		Emoji:      reaction.Emoji,
		CreatedAt:  reaction.CreatedAt,
	}
}

func ReactionDBToEntity(reaction *dbmodel.Reaction) *entitymodel.Reaction {
	if reaction == nil {
		return nil
	}

	entityReaction := &entitymodel.Reaction{
		ID:         reaction.ID,
		SessionID:  reaction.SessionID,
		FromUserID: reaction.FromUserID,
		ToUserID:   reaction.ToUserID,
		Emoji:      reaction.Emoji,
		CreatedAt:  &reaction.CreatedAt,
	}

	return entityReaction
}

func ReactionAPIToEntity(reaction *apimodel.Reaction) *entitymodel.Reaction {
	if reaction == nil {
		return nil
	}

	return &entitymodel.Reaction{
		ID:         reaction.ID,
		SessionID:  reaction.SessionID,
		FromUserID: reaction.FromUserID,
		ToUserID:   reaction.ToUserID,
		Emoji:      reaction.Emoji,
		CreatedAt:  reaction.CreatedAt,
	}
}

func ReactionEntityToDB(reaction *entitymodel.Reaction) *dbmodel.Reaction {
	if reaction == nil {
		return nil
	}

	dbReaction := &dbmodel.Reaction{
		ID:         reaction.ID,
		SessionID:  reaction.SessionID,
		FromUserID: reaction.FromUserID,
		ToUserID:   reaction.ToUserID,
		Emoji:      reaction.Emoji,
	}

	// Конвертируем время
	if reaction.CreatedAt != nil {
		dbReaction.CreatedAt = *reaction.CreatedAt
	}

	return dbReaction
}
