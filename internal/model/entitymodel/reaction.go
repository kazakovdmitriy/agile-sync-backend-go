package entitymodel

import "time"

type Reaction struct {
	ID         string
	SessionID  string
	FromUserID string
	ToUserID   string
	Emoji      string
	CreatedAt  *time.Time
}
