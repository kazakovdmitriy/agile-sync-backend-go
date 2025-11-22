package dbmodel

import "github.com/google/uuid"

type UsersInSessions struct {
	SessionUUID uuid.UUID `db:"session_id"`
	UsesCount   int       `db:"connected_users_count"`
}
