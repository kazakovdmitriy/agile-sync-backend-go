package service

import (
	"backend_go/internal/model"
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/converter"
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SessionService struct {
	sessionRepo SessionRepository
	votesRepo   VoteRepository
	userRepo    UserRepository
	log         *zap.Logger
}

func NewSessionService(
	sessionRepo SessionRepository,
	votesRepo VoteRepository,
	userRepo UserRepository,
	log *zap.Logger,
) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		votesRepo:   votesRepo,
		userRepo:    userRepo,
		log:         log,
	}
}

func (s *SessionService) GetUserSession(
	ctx context.Context,
	userId string,
) ([]apimodel.UserSessions, error) {
	sessions, err := s.sessionRepo.GetByCreator(ctx, userId)
	if err != nil {
		return nil, err
	}

	usersInSession, err := s.sessionRepo.GetCountUsersInUserSessions(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := make([]apimodel.UserSessions, 0)
	for _, session := range sessions {
		usersCount := usersInSession[session.ID]
		apiSession := converter.SessionToUserSession(&session, usersCount)
		response = append(response, *apiSession)
	}

	return response, nil
}

func (s *SessionService) CreateSession(
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

func (s *SessionService) DeleteSession(ctx context.Context, sessionId, userId string) error {
	session, err := s.sessionRepo.GetByID(ctx, sessionId)
	if err != nil {
		return err
	}

	if session.CreatorID.String() != userId {
		return fmt.Errorf("user is not a creator of session")
	}

	return s.sessionRepo.DeleteSession(ctx, sessionId)
}

func (s *SessionService) GetSessionByID(ctx context.Context, sessionId string) (*apimodel.Session, error) {
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

func (s *SessionService) ConnectUser(ctx context.Context, userID, sessionID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return err
	}

	err = s.userRepo.SetOnSession(ctx, userUUID, true)
	if err != nil {
		return err
	}

	return s.sessionRepo.ConnectUser(ctx, userUUID, sessionUUID)
}

func (s *SessionService) DisconnectUser(ctx context.Context, userID, sessionID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return err
	}

	err = s.userRepo.SetOnSession(ctx, userUUID, false)
	if err != nil {
		return err
	}

	return s.sessionRepo.DisconnectUser(ctx, userUUID, sessionUUID)
}

func (s *SessionService) RevealCardsInSession(ctx context.Context, sessionId uuid.UUID, isReveal bool) error {
	return s.sessionRepo.RevealCardsInSession(ctx, sessionId, isReveal)
}

func (s *SessionService) AutoRevealCardsInSession(ctx context.Context, sessionId uuid.UUID) error {
	return s.sessionRepo.AutoRevealCardsInSession(ctx, sessionId)
}
