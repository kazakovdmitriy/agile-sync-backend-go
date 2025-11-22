package service

import (
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/repository"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userService struct {
	userRepository repository.UserRepository
	log            *zap.Logger
}

func NewUserService(userRepository repository.UserRepository, log *zap.Logger) *userService {
	return &userService{
		userRepository: userRepository,
		log:            log,
	}
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*entitymodel.User, error) {
	return s.userRepository.GetByID(ctx, userID)
}
