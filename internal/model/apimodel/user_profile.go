package apimodel

import (
	"backend_go/internal/model/entitymodel"
	"github.com/google/uuid"
)

// UserProfile профиль пользователя
// @Description Полная информация о пользователе
type UserProfile struct {
	Id            uuid.UUID                  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string                     `json:"name" example:"John Doe"`
	Email         *string                    `json:"email" example:"john@example.com"`
	IsActive      bool                       `json:"is_active" example:"true"`
	IsVerified    bool                       `json:"is_verified" example:"false"`
	IsGuest       bool                       `json:"is_guest" example:"false"`
	OAuthProvider *entitymodel.OAuthProvider `json:"oauth_provider,omitempty" swaggertype:"string" example:"google"`
	AvatarUrl     string                     `json:"avatar_url" example:"https://example.com/avatar.jpg"`
}
