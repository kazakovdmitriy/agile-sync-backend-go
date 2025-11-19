package entitymodel

import (
	"github.com/google/uuid"
	"time"
)

type Vote struct {
	ID        uuid.UUID  `db:"id"`
	SessionID uuid.UUID  `db:"session_id"`
	UserID    uuid.UUID  `db:"user_id"`
	Value     string     `db:"value"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
