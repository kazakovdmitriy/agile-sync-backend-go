package service

import (
	"backend_go/internal/mocks"
	"backend_go/internal/model"
	"context"
	"go.uber.org/mock/gomock"
	"testing"

	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSessionService_GetUserSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockVotesRepo := mocks.NewMockVoteRepository(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewSessionService(mockSessionRepo, mockVotesRepo, logger)

	ctx := context.Background()
	userID := uuid.New() // ← валидный UUID
	userIDStr := userID.String()

	sessions := []*entitymodel.Session{
		{
			ID:        uuid.New(),
			Name:      "Session 1",
			CreatorID: userID, // ← uuid.UUID, не string!
		},
	}

	mockSessionRepo.EXPECT().GetByCreator(ctx, userIDStr).Return(sessions, nil)

	result, err := service.GetUserSession(ctx, userIDStr)

	assert.NoError(t, err)
	assert.Equal(t, sessions, result)
}

func TestSessionService_CreateSession_GuestUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockVotesRepo := mocks.NewMockVoteRepository(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewSessionService(mockSessionRepo, mockVotesRepo, logger)

	user := &entitymodel.User{ID: uuid.New(), Name: "Guest User", IsGuest: true}
	sessionCreate := &apimodel.SessionCreate{Name: "Test Session", DeckType: "fibonacci"}

	_, err := service.CreateSession(context.Background(), sessionCreate, user)

	assert.Error(t, err)
	assert.Equal(t, "user is a guest", err.Error())
}

func TestSessionService_CreateSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockVotesRepo := mocks.NewMockVoteRepository(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewSessionService(mockSessionRepo, mockVotesRepo, logger)

	userId := uuid.New()
	user := &entitymodel.User{ID: userId, Name: "Alice", IsGuest: false}
	sessionCreate := &apimodel.SessionCreate{Name: "Planning Session", DeckType: "standard"}

	expectedSession := &entitymodel.Session{
		Name:          sessionCreate.Name,
		DeckType:      sessionCreate.DeckType,
		CardsRevealed: false,
		CreatorID:     user.ID,
		CreatorName:   user.Name,
		CreatedVia:    "web",
	}

	mockSessionRepo.EXPECT().
		CreateSession(gomock.Any(), gomock.AssignableToTypeOf(&entitymodel.Session{})).
		DoAndReturn(func(_ context.Context, s *entitymodel.Session) (*entitymodel.Session, error) {
			// Проверяем поля
			assert.Equal(t, sessionCreate.Name, s.Name)
			assert.Equal(t, sessionCreate.DeckType, s.DeckType)
			assert.Equal(t, user.ID, s.CreatorID)
			assert.Equal(t, user.Name, s.CreatorName)
			assert.Equal(t, "web", s.CreatedVia)
			return expectedSession, nil
		})

	result, err := service.CreateSession(context.Background(), sessionCreate, user)

	assert.NoError(t, err)
	assert.Equal(t, expectedSession, result)
}

func TestSessionService_DeleteSession_NotCreator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockVotesRepo := mocks.NewMockVoteRepository(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewSessionService(mockSessionRepo, mockVotesRepo, logger)

	sessionId := uuid.New().String()
	userId := uuid.New().String()
	creatorId := uuid.New()

	session := &entitymodel.Session{
		ID:        uuid.MustParse(sessionId),
		CreatorID: creatorId,
	}

	mockSessionRepo.EXPECT().GetByID(gomock.Any(), sessionId).Return(session, nil)

	err := service.DeleteSession(context.Background(), sessionId, userId)

	assert.Error(t, err)
	assert.Equal(t, "user is not a creator of session", err.Error())
}

func TestSessionService_DeleteSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockVotesRepo := mocks.NewMockVoteRepository(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewSessionService(mockSessionRepo, mockVotesRepo, logger)

	sessionId := uuid.New().String()
	userId := uuid.New().String()

	session := &entitymodel.Session{
		ID:        uuid.MustParse(sessionId),
		CreatorID: uuid.MustParse(userId),
	}

	mockSessionRepo.EXPECT().GetByID(gomock.Any(), sessionId).Return(session, nil)
	mockSessionRepo.EXPECT().DeleteSession(gomock.Any(), sessionId).Return(nil)

	err := service.DeleteSession(context.Background(), sessionId, userId)

	assert.NoError(t, err)
}

func TestSessionService_GetSessionByID_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockVotesRepo := mocks.NewMockVoteRepository(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewSessionService(mockSessionRepo, mockVotesRepo, logger)

	_, err := service.GetSessionByID(context.Background(), "not-a-uuid")

	assert.Error(t, err)
}

func TestSessionService_GetSessionByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockVotesRepo := mocks.NewMockVoteRepository(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewSessionService(mockSessionRepo, mockVotesRepo, logger)

	sessionID := uuid.New()
	sessionIDStr := sessionID.String()

	sessionDB := &entitymodel.Session{
		ID:            sessionID,
		Name:          "Estimation",
		DeckType:      "fibonacci",
		CardsRevealed: true,
		CreatorID:     uuid.New(),
		CreatorName:   "Bob",
		CreatedVia:    "web",
	}

	// Используем apimodel.UsersInSession, как ожидает репозиторий
	users := []apimodel.UsersInSession{
		{ID: uuid.New(), Name: "User1"},
		{ID: uuid.New(), Name: "User2"},
	}

	votes := []*entitymodel.Vote{
		{UserID: users[0].ID, Value: "5"},
		{UserID: users[1].ID, Value: "8"},
	}

	// Моки
	mockSessionRepo.EXPECT().GetByID(gomock.Any(), sessionIDStr).Return(sessionDB, nil)
	mockSessionRepo.EXPECT().GetUsers(gomock.Any(), sessionID).Return(users, nil) // ← правильный тип!
	mockSessionRepo.EXPECT().GetBySessionsID(gomock.Any(), sessionID).Return(votes, nil)

	result, err := service.GetSessionByID(context.Background(), sessionIDStr)

	assert.NoError(t, err)
	assert.Equal(t, sessionDB.ID, result.ID)
	assert.Equal(t, sessionDB.Name, result.Name)
	assert.Equal(t, sessionDB.DeckType, result.DeckType)
	assert.Equal(t, sessionDB.CardsRevealed, result.CardsRevealed)
	assert.Equal(t, sessionDB.CreatorID, result.CreatorID)
	assert.Equal(t, sessionDB.CreatorName, result.CreatorName)
	assert.Equal(t, model.DeckValues["fibonacci"], result.DeckValues)
	assert.Equal(t, users, result.Users) // apimodel.UsersInSession

	expectedVotes := map[uuid.UUID]string{
		users[0].ID: "5",
		users[1].ID: "8",
	}
	assert.Equal(t, expectedVotes, result.Votes)
}

func TestSessionService_ConnectUserToSession_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockVotesRepo := mocks.NewMockVoteRepository(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewSessionService(mockSessionRepo, mockVotesRepo, logger)

	err := service.ConnectUserToSession(context.Background(), "bad-uuid", "also-bad")

	assert.Error(t, err)
}

func TestSessionService_ConnectUserToSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockVotesRepo := mocks.NewMockVoteRepository(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewSessionService(mockSessionRepo, mockVotesRepo, logger)

	userID := uuid.New()
	sessionID := uuid.New()

	mockSessionRepo.EXPECT().
		ConnectUserToSession(gomock.Any(), userID, sessionID).
		Return(nil)

	err := service.ConnectUserToSession(context.Background(), userID.String(), sessionID.String())

	assert.NoError(t, err)
}
