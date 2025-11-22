package service

import (
	"backend_go/internal/model/entitymodel"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VoteService struct {
	voteRepository VoteRepository
	log            *zap.Logger
}

func NewVoteService(voteRepository VoteRepository, log *zap.Logger) *VoteService {
	return &VoteService{
		voteRepository: voteRepository,
		log:            log,
	}
}

func (s *VoteService) SaveVote(ctx context.Context, vote *entitymodel.Vote) (uuid.UUID, error) {
	return s.voteRepository.SetVoteValue(ctx, vote.SessionID, vote.UserID, vote.Value)
}

func (s *VoteService) DeleteVoteInSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.voteRepository.DeleteVoteInSession(ctx, sessionID)
}
