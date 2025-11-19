package service

import (
	"backend_go/internal/infrastructure/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

// helper для создания тестовой конфигурации с TTL
func testJWTConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			SecretKey:       "very-secret-key-for-testing-only-32-bytes!",
			Issuer:          "test-issuer",
			AccessTokenTTL:  15,
			RefreshTokenTTL: 7,
		},
	}
}

func TestJwtService_GenerateTokenPair_Success(t *testing.T) {
	cfg := testJWTConfig()
	logger, _ := zap.NewDevelopment()
	svc := NewJwtService(cfg, logger)

	userID := "user-123"
	tokens, err := svc.GenerateTokenPair(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokens["access_token"])
	assert.NotEmpty(t, tokens["refresh_token"])
	assert.Equal(t, "Bearer", tokens["token_type"])

	// Проверяем access token
	parsedAT, err := jwt.ParseWithClaims(tokens["access_token"], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.SecretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedAT.Valid)

	if claims, ok := parsedAT.Claims.(*CustomClaims); ok {
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, userID, claims.Subject)
		assert.Equal(t, cfg.JWT.Issuer, claims.Issuer)
		assert.True(t, claims.ExpiresAt.After(time.Now()), "access token должен быть не просрочен")
	} else {
		assert.Fail(t, "claims не приводятся к *CustomClaims")
	}

	// Проверяем refresh token
	parsedRT, err := jwt.ParseWithClaims(tokens["refresh_token"], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.SecretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedRT.Valid)

	if claims, ok := parsedRT.Claims.(*CustomClaims); ok {
		assert.Equal(t, userID, claims.UserID)
		assert.True(t, claims.ExpiresAt.After(parsedAT.Claims.(*CustomClaims).ExpiresAt.Time), "refresh token должен жить дольше access")
	} else {
		assert.Fail(t, "claims не приводятся к *CustomClaims")
	}
}

func TestJwtService_ValidateToken_InvalidSignature(t *testing.T) {
	cfg := testJWTConfig()
	logger, _ := zap.NewDevelopment()
	svc := NewJwtService(cfg, logger)

	// Генерируем токен с другим ключом
	wrongCfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:       "another-secret-key-xxxxxxxxxxxxxxx",
			Issuer:          "test-issuer",
			AccessTokenTTL:  15,
			RefreshTokenTTL: 7,
		},
	}
	wrongSvc := NewJwtService(wrongCfg, logger)
	tokens, _ := wrongSvc.GenerateTokenPair("user-1")

	_, err := svc.ValidateToken(tokens["access_token"])
	assert.Error(t, err)
	assert.ErrorIs(t, err, jwt.ErrSignatureInvalid)
}

func TestJwtService_ValidateToken_Expired(t *testing.T) {
	cfg := testJWTConfig()
	logger, _ := zap.NewDevelopment()
	svc := NewJwtService(cfg, logger)

	// Создаём вручную просроченный токен
	expiredClaims := CustomClaims{
		UserID: "user-1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-20 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-20 * time.Minute)),
			Issuer:    cfg.JWT.Issuer,
			Subject:   "user-1",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, _ := token.SignedString([]byte(cfg.JWT.SecretKey))

	_, err := svc.ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJwtService_ExtractUserIDFromToken_Success(t *testing.T) {
	cfg := testJWTConfig()
	logger, _ := zap.NewDevelopment()
	svc := NewJwtService(cfg, logger)

	tokens, err := svc.GenerateTokenPair("user-789")
	assert.NoError(t, err)

	userID, err := svc.ExtractUserIDFromToken(tokens["access_token"])
	assert.NoError(t, err)
	assert.Equal(t, "user-789", userID)
}

func TestJwtService_ExtractUserIDFromToken_InvalidToken(t *testing.T) {
	cfg := testJWTConfig()
	logger, _ := zap.NewDevelopment()
	svc := NewJwtService(cfg, logger)

	userID, err := svc.ExtractUserIDFromToken("invalid.token.string")
	assert.Error(t, err)
	assert.Empty(t, userID)
}

func TestJwtService_RefreshToken_Success(t *testing.T) {
	cfg := testJWTConfig()
	logger, _ := zap.NewDevelopment()
	svc := NewJwtService(cfg, logger)

	tokens, err := svc.GenerateTokenPair("user-refresh")
	assert.NoError(t, err)

	newTokens, err := svc.RefreshToken(tokens["refresh_token"])
	assert.NoError(t, err)
	assert.NotEmpty(t, newTokens["access_token"])
	assert.NotEmpty(t, newTokens["refresh_token"])
	assert.Equal(t, "Bearer", newTokens["token_type"])

	// Проверяем, что новый access token содержит правильный user ID
	userID, err := svc.ExtractUserIDFromToken(newTokens["access_token"])
	assert.NoError(t, err)
	assert.Equal(t, "user-refresh", userID)
}

func TestJwtService_RefreshToken_WithAccessToken_ShouldFail(t *testing.T) {
	cfg := testJWTConfig()
	logger, _ := zap.NewDevelopment()
	svc := NewJwtService(cfg, logger)

	tokens, err := svc.GenerateTokenPair("user-abc")
	assert.NoError(t, err)

	_, err = svc.RefreshToken(tokens["access_token"]) // передаём access вместо refresh

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token type")
}

func TestJwtService_RefreshToken_InvalidToken(t *testing.T) {
	cfg := testJWTConfig()
	logger, _ := zap.NewDevelopment()
	svc := NewJwtService(cfg, logger)

	_, err := svc.RefreshToken("totally.fake.token")
	assert.Error(t, err)
}
