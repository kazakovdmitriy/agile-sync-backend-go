package service

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type sessionService struct {
	sessionRepo repository.SessionRepository
	votesRepo   repository.VoteRepository
	log         *zap.Logger
}

func NewSessionService(
	sessionRepo repository.SessionRepository,
	votesRepo repository.VoteRepository,
	log *zap.Logger,
) *sessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
		votesRepo:   votesRepo,
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

	err = s.sessionRepo.ConnectUserToSession(ctx, user.ID, sessionResult.ID)
	if err != nil {
		s.log.Error("Failed to connect user to session", zap.Error(err))
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

func (s *sessionService) GetSessionByID(ctx context.Context, sessionId string) (*apimodel.Session, error) {
	sessionUUID, err := uuid.Parse(sessionId)
	if err != nil {
		return nil, err
	}

	votes, err := s.votesRepo.GetBySessionsID(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse("623b7dea-b040-4f9c-9eb7-357600a85f4f")
	if err != nil {
		return nil, err
	}

	for range 5 {
		err := s.votesRepo.SetVoteValue(ctx, sessionUUID, userUUID, "8")
		if err != nil {
			return nil, err
		}
		time.Sleep(60 * time.Millisecond)
	}

	fmt.Println(votes)

	return nil, nil
}
