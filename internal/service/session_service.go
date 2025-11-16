package service

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/repository"
	"context"
	"fmt"
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

func (s *sessionService) CreateSession(
	ctx context.Context,
	sessionCreate *apimodel.SessionCreate,
	user *entitymodel.User,
) (*entitymodel.Session, error) {
	if user.IsGuest {
		return nil, fmt.Errorf("user is a guest")
	}

	session := entitymodel.Session{
		Name:          sessionCreate.Name,
		DeckType:      sessionCreate.DeckType,
		CardsRevealed: false,
		CreatorID:     user.ID,
		CreatorName:   user.Name,
		CreatedVia:    "web",
	}

	sessionResult, err := s.sessionRepo.CreateSession(ctx, &session)
	if err != nil {
		return nil, err
	}

	return sessionResult, nil
}

func (s *sessionService) DeleteSession(ctx context.Context, sessionId string, userId string) error {

	session, err := s.sessionRepo.GetByID(ctx, sessionId)
	if err != nil {
		return err
	}

	if session.CreatorID.String() != userId {
		return fmt.Errorf("user is not a creator of session")
	}

	return s.sessionRepo.DeleteSession(ctx, sessionId)
}
