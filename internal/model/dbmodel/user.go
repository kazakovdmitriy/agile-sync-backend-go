package dbmodel

import "time"

type OAuthProviderEnum string

const (
	Google OAuthProviderEnum = "google"
	Yandex OAuthProviderEnum = "yandex"
)

type User struct {
	ID             string             `db:"id"`
	Name           string             `db:"name"`
	Email          string             `db:"email"`
	HashedPassword string             `db:"hashed_password"`
	IsActive       bool               `db:"is_active"`
	IsVerified     bool               `db:"is_verified"`
	IsGuest        bool               `db:"is_guest"`
	OAuthProvider  *OAuthProviderEnum `db:"oauth_provider"`
	OAuthID        *string            `db:"oauth_id"`
	AvatarURL      *string            `db:"avatar_url"`
	IsCreator      bool               `db:"is_creator"`
	IsWatcher      bool               `db:"is_watcher"`
	OnSession      bool               `db:"on_session"`
	SocketID       *string            `db:"socket_id"`
	CreatedAt      time.Time          `db:"created_at"`
	UpdatedAt      *time.Time         `db:"updated_at"`
}
