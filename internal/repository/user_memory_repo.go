package repository

import (
	"backend_go/internal/model/entitymodel"
	"context"
	"go.uber.org/zap"
)

type UserMemoryRepo struct {
	users []entitymodel.User
	log   *zap.Logger
}

func NewUserMemoryRepo(log *zap.Logger) *UserMemoryRepo {
	return &UserMemoryRepo{
		log: log,
	}
}

func (repo *UserMemoryRepo) Create(ctx context.Context, user entitymodel.User) {
	repo.users = append(repo.users, user)
}

func (repo *UserMemoryRepo) GetByEmail(ctx context.Context, email string) (*entitymodel.User, error) {
	for _, user := range repo.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, nil
}
