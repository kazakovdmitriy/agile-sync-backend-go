package service

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req *apimodel.UserRegister) (*apimodel.TokenResponse, error)
	GuestLogin(ctx context.Context, req *apimodel.GuestLogin) (*apimodel.TokenResponse, error)
	Login(ctx context.Context, req *apimodel.UserLogin) (*apimodel.TokenResponse, error)
	ValidateToken(ctx context.Context, token string) (*entitymodel.User, error)
}

type JWTService interface {
	GenerateTokenPair(userID string) (map[string]string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	ExtractUserIDFromToken(tokenString string) (string, error)
	RefreshToken(refreshToken string) (map[string]string, error)
}

type SessionService interface {
	GetUserSession(ctx context.Context, userId string) ([]*entitymodel.Session, error)
	CreateSession(ctx context.Context, sessions *apimodel.SessionCreate, user *entitymodel.User) (*entitymodel.Session, error)
	DeleteSession(ctx context.Context, sessionId, userId string) error
	GetSessionByID(ctx context.Context, sessionId string) (*apimodel.Session, error)
	ConnectUser(ctx context.Context, userID, sessionID string) error
	DisconnectUser(ctx context.Context, userID, sessionID string) error
	RevealCardsInSession(ctx context.Context, sessionId uuid.UUID, isReveal bool) error
}

type UserService interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*entitymodel.User, error)
}

type VoteService interface {
	SaveVote(ctx context.Context, vote *entitymodel.Vote) (uuid.UUID, error)
	DeleteVoteInSession(ctx context.Context, sessionID uuid.UUID) error
}
