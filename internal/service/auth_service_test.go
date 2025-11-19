package service

import (
	"backend_go/internal/mocks"
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
	"backend_go/pkg/hash"
	"context"
	"database/sql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

func TestAuthService_Register_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewAuthService(mockUserRepo, mockJWTService, logger)

	req := &apimodel.UserRegister{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	}

	existingUser := &entitymodel.User{
		ID:    uuid.New(),
		Name:  "Alice",
		Email: &req.Email,
	}

	mockUserRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(existingUser, nil)

	_, err := service.Register(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUserAlreadyExists)
	assert.Contains(t, err.Error(), req.Email)
}

func TestAuthService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewAuthService(mockUserRepo, mockJWTService, logger)

	req := &apimodel.UserRegister{
		Name:     "Charlie",
		Email:    "charlie@example.com",
		Password: "securePass!123",
	}

	userID := uuid.New()
	createdUser := &entitymodel.User{
		ID:             userID,
		Name:           req.Name,
		Email:          &req.Email,
		HashedPassword: "hashed_password_stub",
		IsActive:       true,
		IsVerified:     false,
		IsGuest:        false,
		IsCreator:      false,
		IsWatcher:      false,
		OnSession:      true,
	}

	tokens := map[string]string{
		"access_token":  "fake-access-token",
		"refresh_token": "fake-refresh-token",
		"token_type":    "Bearer",
	}

	mockUserRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(nil, sql.ErrNoRows)
	mockUserRepo.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(&entitymodel.User{})).
		Return(createdUser, nil)
	mockJWTService.EXPECT().
		GenerateTokenPair(userID.String()).
		Return(tokens, nil)

	resp, err := service.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, tokens["access_token"], resp.AccessToken)
	assert.Equal(t, tokens["refresh_token"], resp.RefreshToken)
	assert.Equal(t, tokens["token_type"], resp.TokenType)
}

func TestAuthService_GuestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewAuthService(mockUserRepo, mockJWTService, logger)

	req := &apimodel.GuestLogin{Name: "GuestUser"}

	userID := uuid.New()
	createdUser := &entitymodel.User{
		ID:         userID,
		Name:       req.Name,
		Email:      nil,
		IsActive:   true,
		IsVerified: false,
		IsGuest:    true,
		OnSession:  true,
	}

	tokens := map[string]string{
		"access_token":  "guest-access",
		"refresh_token": "guest-refresh",
		"token_type":    "Bearer",
	}

	mockUserRepo.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(&entitymodel.User{})).
		Return(createdUser, nil)
	mockJWTService.EXPECT().
		GenerateTokenPair(userID.String()).
		Return(tokens, nil)

	resp, err := service.GuestLogin(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, tokens["access_token"], resp.AccessToken)
	assert.Equal(t, tokens["refresh_token"], resp.RefreshToken)
	assert.Equal(t, tokens["token_type"], resp.TokenType)
}

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewAuthService(mockUserRepo, mockJWTService, logger)

	req := &apimodel.UserLogin{
		Email:    "active@example.com",
		Password: "mySecretPass",
	}

	hashedPass, _ := hash.HashPassword(req.Password)
	userID := uuid.New()
	user := &entitymodel.User{
		ID:             userID,
		Email:          &req.Email,
		HashedPassword: hashedPass,
		IsActive:       true,
	}

	tokens := map[string]string{
		"access_token":  "access-123",
		"refresh_token": "refresh-123",
		"token_type":    "Bearer",
	}

	mockUserRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(user, nil)
	mockJWTService.EXPECT().GenerateTokenPair(userID.String()).Return(tokens, nil)

	resp, err := service.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, tokens["access_token"], resp.AccessToken)
	assert.Equal(t, tokens["refresh_token"], resp.RefreshToken)
	assert.Equal(t, tokens["token_type"], resp.TokenType)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewAuthService(mockUserRepo, mockJWTService, logger)

	req := &apimodel.UserLogin{
		Email:    "user@example.com",
		Password: "wrongpass",
	}

	user := &entitymodel.User{
		ID:             uuid.New(),
		Email:          &req.Email,
		HashedPassword: "$2a$10$validhash...", // но не от "wrongpass"
		IsActive:       true,
	}

	mockUserRepo.EXPECT().GetByEmail(gomock.Any(), req.Email).Return(user, nil)

	_, err := service.Login(context.Background(), req)

	assert.Error(t, err)
	// Обычно это bcrypt.ErrMismatchedHashAndPassword, но точный тип зависит от pkg/hash
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewAuthService(mockUserRepo, mockJWTService, logger)

	userID := uuid.New()
	token := "valid.jwt.token"

	mockJWTService.EXPECT().ExtractUserIDFromToken(token).Return(userID.String(), nil)
	mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(&entitymodel.User{ID: userID}, nil)

	user, err := service.ValidateToken(context.Background(), token)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
}

func TestAuthService_ValidateToken_ExtractFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	logger, _ := zap.NewDevelopment()

	service := NewAuthService(mockUserRepo, mockJWTService, logger)

	token := "invalid.token"

	mockJWTService.EXPECT().ExtractUserIDFromToken(token).Return("", jwt.ErrSignatureInvalid)

	_, err := service.ValidateToken(context.Background(), token)

	assert.Error(t, err)
	assert.ErrorIs(t, err, jwt.ErrSignatureInvalid)
}
