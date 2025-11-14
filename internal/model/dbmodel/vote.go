package dbmodel

import "time"

type Vote struct {
	ID        string     `db:"id"`
	SessionID string     `db:"session_id"`
	UserID    string     `db:"user_id"`
	Value     string     `db:"value"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
