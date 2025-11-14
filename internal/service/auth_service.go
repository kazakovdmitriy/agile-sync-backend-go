package service

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/repository"
	"backend_go/pkg/hash"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthServiceImpl struct {
	userRepo   repository.UserRepository
	jwtService JWTService
	log        *zap.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwtService JWTService, log *zap.Logger) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
		log:        log,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *apimodel.UserRegister) (*apimodel.TokenResponse, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := entitymodel.User{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		IsActive:       true,
		IsVerified:     false,
		IsGuest:        false,
		IsCreator:      false,
		IsWatcher:      false,
		OnSession:      true,
	}

	createdUser, err := s.userRepo.Create(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwtService.GenerateTokenPair(createdUser.ID.String(), createdUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &apimodel.TokenResponse{
		AccessToken:  tokens["access_token"],
		RefreshToken: tokens["refresh_token"],
		TokenType:    tokens["token_type"],
	}, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *apimodel.UserLogin) (*apimodel.TokenResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.log.Info("failed to find user by email", zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}

	err = hash.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		s.log.Info("failed to check password", zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}

	if user.IsActive != true {
		s.log.Info("user is not active", zap.String("email", req.Email))
		return nil, fmt.Errorf("user is not active")
	}

	tokens, err := s.jwtService.GenerateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &apimodel.TokenResponse{
		AccessToken:  tokens["access_token"],
		RefreshToken: tokens["refresh_token"],
		TokenType:    tokens["token_type"],
	}, nil
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (*entitymodel.User, error) {
	userIDStr, err := s.jwtService.ExtractUserIDFromToken(token)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.log.Info("failed to parse userID from token", zap.String("token", token))
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
