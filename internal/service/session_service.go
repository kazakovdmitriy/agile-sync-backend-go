package service

import (
	"backend_go/internal/model"
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/repository"
	"backend_go/internal/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

	return sessionResult, nil
}

func (s *sessionService) DeleteSession(ctx context.Context, sessionId, userId string) error {
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

	sessionDB, err := s.sessionRepo.GetByID(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	users, err := s.sessionRepo.GetUsers(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	votes, err := s.votesRepo.GetVotesInSessions(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	// Собираем настоящие голоса (user_id → значение)
	realUserVotes := make(map[uuid.UUID]string)
	for _, vote := range votes {
		realUserVotes[vote.UserID] = vote.Value
	}

	// Определяем, какие голоса отдавать клиенту
	var clientVotes map[uuid.UUID]string
	if sessionDB.CardsRevealed {
		clientVotes = realUserVotes
	} else {
		// Карты скрыты — отдаём "hidden" для каждого проголосовавшего
		clientVotes = make(map[uuid.UUID]string)
		for userID := range realUserVotes {
			clientVotes[userID] = "hidden"
		}
	}

	userVotes := make(map[uuid.UUID]string)
	for _, vote := range votes {
		userVotes[vote.UserID] = vote.Value
	}

	var session = &apimodel.Session{
		ID:            sessionDB.ID,
		Name:          sessionDB.Name,
		DeckType:      sessionDB.DeckType,
		CardsRevealed: sessionDB.CardsRevealed,
		CreatorID:     sessionDB.CreatorID,
		CreatorName:   sessionDB.CreatorName,
		AutoReveal:    sessionDB.AutoReveal,
		AllowEmoji:    sessionDB.AllowEmoji,
		CreatedVia:    sessionDB.CreatedVia,
		DeckValues:    model.DeckValues[sessionDB.DeckType],
		Users:         users,
		Votes:         clientVotes,
	}

	if sessionDB.CardsRevealed {
		session.MostCommonVote = utils.FindMostPopularVote(votes, sessionDB.DeckType)
	}

	return session, nil
}

func (s *sessionService) ConnectUserToSession(ctx context.Context, userID, sessionID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return err
	}

	return s.sessionRepo.ConnectUserToSession(ctx, userUUID, sessionUUID)
}

func (s *sessionService) RevealCardsInSession(ctx context.Context, sessionId uuid.UUID, isReveal bool) error {
	return s.sessionRepo.RevealCardsInSession(ctx, sessionId, isReveal)
}
