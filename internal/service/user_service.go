package service

import (
	"backend_go/internal/model/entitymodel"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepository UserRepository
	log            *zap.Logger
}

func NewUserService(userRepository UserRepository, log *zap.Logger) *UserService {
	return &UserService{
		userRepository: userRepository,
		log:            log,
	}
}

func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (*entitymodel.User, error) {
	return s.userRepository.GetByID(ctx, userID)
}
