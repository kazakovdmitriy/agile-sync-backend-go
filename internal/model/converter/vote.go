package converter

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/dbmodel"
	"backend_go/internal/model/entitymodel"
)

func VoteEntityToAPI(vote *entitymodel.Vote) *apimodel.Vote {
	if vote == nil {
		return nil
	}

	return &apimodel.Vote{
		ID:        vote.ID,
		SessionID: vote.SessionID,
		UserID:    vote.UserID,
		Value:     vote.Value,
		CreatedAt: vote.CreatedAt,
		UpdatedAt: vote.UpdatedAt,
	}
}

func VoteDBToEntity(vote *dbmodel.Vote) *entitymodel.Vote {
	if vote == nil {
		return nil
	}

	entityVote := &entitymodel.Vote{
		ID:        vote.ID,
		SessionID: vote.SessionID,
		UserID:    vote.UserID,
		Value:     vote.Value,
		CreatedAt: &vote.CreatedAt,
	}

	// Конвертируем UpdatedAt
	if vote.UpdatedAt != nil {
		updatedAt := *vote.UpdatedAt
		entityVote.UpdatedAt = &updatedAt
	}

	return entityVote
}

func VoteAPIToEntity(vote *apimodel.Vote) *entitymodel.Vote {
	if vote == nil {
		return nil
	}

	return &entitymodel.Vote{
		ID:        vote.ID,
		SessionID: vote.SessionID,
		UserID:    vote.UserID,
		Value:     vote.Value,
		CreatedAt: vote.CreatedAt,
		UpdatedAt: vote.UpdatedAt,
	}
}

func VoteEntityToDB(vote *entitymodel.Vote) *dbmodel.Vote {
	if vote == nil {
		return nil
	}

	dbVote := &dbmodel.Vote{
		ID:        vote.ID,
		SessionID: vote.SessionID,
		UserID:    vote.UserID,
		Value:     vote.Value,
	}

	// Конвертируем время
	if vote.CreatedAt != nil {
		dbVote.CreatedAt = *vote.CreatedAt
	}
	if vote.UpdatedAt != nil {
		updatedAt := *vote.UpdatedAt
		dbVote.UpdatedAt = &updatedAt
	}

	return dbVote
}
