package service

import (
	"backend_go/internal/infrastructure/config"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"time"
)

type jwtService struct {
	cfg *config.Config
	log *zap.Logger
}

func NewJwtService(cfg *config.Config, log *zap.Logger) *jwtService {
	return &jwtService{
		cfg: cfg,
		log: log,
	}
}

// CustomClaims - кастомные claims для нашего приложения
type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *jwtService) GenerateTokenPair(userID string, email string) (map[string]string, error) {
	// Access Token
	accessTokenClaims := CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.GetAccessTokenTTL())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.cfg.JWT.Issuer,
			Subject:   userID,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWT.SecretKey))
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshTokenClaims := CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.GetRefreshTokenTTL())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.cfg.JWT.Issuer,
			Subject:   userID,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWT.SecretKey))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
		"token_type":    "Bearer",
	}, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.cfg.JWT.SecretKey), nil
	})
}

func (s *jwtService) ExtractUserIDFromToken(tokenString string) (string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", jwt.ErrInvalidKey
}

func (s *jwtService) RefreshToken(refreshToken string) (map[string]string, error) {
	token, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return s.GenerateTokenPair(claims.UserID, claims.Email)
}
