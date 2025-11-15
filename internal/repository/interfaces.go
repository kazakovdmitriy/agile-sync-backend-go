package repository

import (
	"backend_go/internal/model/entitymodel"
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entitymodel.User) (*entitymodel.User, error)
	GetByEmail(ctx context.Context, email string) (*entitymodel.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entitymodel.User, error)
}

type SessionRepository interface {
	GetByCreator(ctx context.Context, userId string) ([]*entitymodel.Session, error)
	CreateSession(ctx context.Context, session *entitymodel.Session) (*entitymodel.Session, error)
}
