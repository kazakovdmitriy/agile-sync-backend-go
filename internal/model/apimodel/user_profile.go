package apimodel

import (
	"backend_go/internal/model/entitymodel"
	"github.com/google/uuid"
)

type UserProfile struct {
	Id            uuid.UUID                  `json:"id"`
	Name          string                     `json:"name"`
	Email         *string                    `json:"email"`
	IsActive      bool                       `json:"is_active"`
	IsVerified    bool                       `json:"is_verified"`
	IsGuest       bool                       `json:"is_guest"`
	OAuthProvider *entitymodel.OAuthProvider `json:"oauth_provider,omitempty"`
	AvatarUrl     string                     `json:"avatar_url"`
}
