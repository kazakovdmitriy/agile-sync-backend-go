package entitymodel

import "time"

type Vote struct {
	ID        string
	SessionID string
	UserID    string
	Value     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
