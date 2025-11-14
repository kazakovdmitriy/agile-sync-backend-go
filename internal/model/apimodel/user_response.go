package apimodel

import (
	"time"
)

type OAuthProvider string

const (
	Google OAuthProvider = "google"
	Yandex OAuthProvider = "yandex"
)

type UserResponse struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email,omitempty"`
	HashedPassword string         `json:"-"`
	IsActive       bool           `json:"is_active"`
	IsVerified     bool           `json:"is_verified"`
	IsGuest        bool           `json:"is_guest"`
	OAuthProvider  *OAuthProvider `json:"oauth_provider,omitempty"`
	OAuthID        *string        `json:"oauth_id,omitempty"`
	AvatarURL      *string        `json:"avatar_url,omitempty"`
	SessionID      *string        `json:"session_id,omitempty"`
	IsCreator      bool           `json:"is_creator"`
	IsWatcher      bool           `json:"is_watcher"`
	OnSession      bool           `json:"on_session"`
	SocketID       *string        `json:"socket_id,omitempty"`
	CreatedAt      *time.Time     `json:"created_at,omitempty"`
	UpdatedAt      *time.Time     `json:"updated_at,omitempty"`
}

func (u *UserResponse) IsOAuthUser() bool {
	return u.OAuthProvider != nil && u.OAuthID != nil
}

func (u *UserResponse) IsGuestUser() bool {
	return u.IsGuest
}
