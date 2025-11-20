package service

import (
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/repository"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type voteService struct {
	voteRepository repository.VoteRepository
	log            *zap.Logger
}

func NewVoteService(voteRepository repository.VoteRepository, log *zap.Logger) *voteService {
	return &voteService{
		voteRepository: voteRepository,
		log:            log,
	}
}

func (s *voteService) SaveVote(ctx context.Context, vote *entitymodel.Vote) (uuid.UUID, error) {
	return s.voteRepository.SetVoteValue(ctx, vote.SessionID, vote.UserID, vote.Value)
}
