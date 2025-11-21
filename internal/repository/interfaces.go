package repository

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entitymodel.User) (*entitymodel.User, error)
	GetByEmail(ctx context.Context, email string) (*entitymodel.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entitymodel.User, error)
	DeleteInactiveGuests(ctx context.Context, duration string) (int64, error)
}

type SessionRepository interface {
	GetByCreator(ctx context.Context, userId string) ([]*entitymodel.Session, error)
	GetByID(ctx context.Context, sessionId string) (*entitymodel.Session, error)
	CreateSession(ctx context.Context, session *entitymodel.Session) (*entitymodel.Session, error)
	DeleteSession(ctx context.Context, sessionId string) error
	GetUsers(ctx context.Context, sessionID uuid.UUID) ([]apimodel.UsersInSession, error)
	GetBySessionsID(ctx context.Context, sessionID uuid.UUID) ([]*entitymodel.Vote, error)
	ConnectUserToSession(ctx context.Context, userID, sessionID uuid.UUID) error
	DisconnectUserFromSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	RevealCardsInSession(ctx context.Context, sessionID uuid.UUID, isReveal bool) error
}

type VoteRepository interface {
	SetVoteValue(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, vote string) (uuid.UUID, error)
}
