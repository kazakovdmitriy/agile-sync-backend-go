package entitymodel

import (
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
	"time"
)

type OAuthProvider string

const (
	Google OAuthProvider = "google"
	Yandex OAuthProvider = "yandex"
)

type User struct {
	ID             uuid.UUID
	Name           string
	Email          string
	HashedPassword string
	IsActive       bool
	IsVerified     bool
	IsGuest        bool
	OAuthProvider  *OAuthProvider
	OAuthID        *string
	AvatarURL      *string
	IsCreator      bool
	IsWatcher      bool
	OnSession      bool
	SocketID       *string
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

// MarshalLogObject реализует zapcore.ObjectMarshaler для структурированного логирования.
func (u *User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", u.ID.String())
	enc.AddString("name", u.Name)
	enc.AddString("email", u.Email)
	enc.AddBool("is_active", u.IsActive)
	enc.AddBool("is_verified", u.IsVerified)
	enc.AddBool("is_guest", u.IsGuest)
	enc.AddBool("is_creator", u.IsCreator)
	enc.AddBool("is_watcher", u.IsWatcher)
	enc.AddBool("on_session", u.OnSession)

	// Поля-указатели
	if u.HashedPassword != "" {
		enc.AddString("hashed_password", "***")
	}
	if u.OAuthProvider != nil {
		enc.AddString("oauth_provider", string(*u.OAuthProvider))
	}
	if u.OAuthID != nil {
		enc.AddString("oauth_id", *u.OAuthID)
	}
	if u.AvatarURL != nil {
		enc.AddString("avatar_url", *u.AvatarURL)
	}
	if u.SocketID != nil {
		enc.AddString("socket_id", *u.SocketID)
	}
	if u.CreatedAt != nil {
		enc.AddTime("created_at", *u.CreatedAt)
	}
	if u.UpdatedAt != nil {
		enc.AddTime("updated_at", *u.UpdatedAt)
	}

	return nil
}

func (u *User) IsOAuthUser() bool {
	return u.OAuthProvider != nil && u.OAuthID != nil
}

func (u *User) IsGuestUser() bool {
	return u.IsGuest
}
