package service

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	Register(ctx context.Context, req *apimodel.UserRegister) (*apimodel.TokenResponse, error)
	Login(ctx context.Context, req *apimodel.UserLogin) (*apimodel.TokenResponse, error)
	ValidateToken(ctx context.Context, token string) (*entitymodel.User, error)
	//Me(ctx context.Context) (*apimodel.UserResponse, error)
}

type JWTService interface {
	GenerateTokenPair(userID string, email string) (map[string]string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	ExtractUserIDFromToken(tokenString string) (string, error)
	RefreshToken(refreshToken string) (map[string]string, error)
}
