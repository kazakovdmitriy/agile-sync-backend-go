package dbmodel

import "time"

type Reaction struct {
	ID         string    `db:"id"`
	SessionID  string    `db:"session_id"`
	FromUserID string    `db:"from_user_id"`
	ToUserID   string    `db:"to_user_id"`
	Emoji      string    `db:"emoji"`
	CreatedAt  time.Time `db:"created_at"`
}
