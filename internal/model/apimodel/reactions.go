package apimodel

import "time"

type Reaction struct {
	ID         string     `json:"id"`
	SessionID  string     `json:"session_id"`
	FromUserID string     `json:"from_user_id"`
	ToUserID   string     `json:"to_user_id"`
	Emoji      string     `json:"emoji"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}
