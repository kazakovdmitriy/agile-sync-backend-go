package service

import (
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/repository"
	"context"
	"go.uber.org/zap"
)

type sessionService struct {
	sessionRepo repository.SessionRepository
	log         *zap.Logger
}

func NewSessionService(
	sessionRepo repository.SessionRepository,
	log *zap.Logger,
) *sessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
		log:         log,
	}
}

func (s *sessionService) GetUserSession(
	ctx context.Context,
	userId string,
) ([]*entitymodel.Session, error) {
	sessions, err := s.sessionRepo.GetByCreator(ctx, userId)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}
