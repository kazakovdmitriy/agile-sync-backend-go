package websocketmodel

import (
	"backend_go/internal/model/entitymodel"
	"github.com/google/uuid"
)

// BaseMessage базовое сообщение WebSocket
type BaseMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// JoinSessionData данные для присоединения к сессии
type JoinSessionData struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	IsWatcher *bool     `json:"is_watcher"`
}

// RevealCardsData данные открытия карт в сессии
type RevealCardsData struct {
	SessionID uuid.UUID `json:"session_id"`
}

// VoteData данные для голосования
type VoteData struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	Value     string    `json:"value"`
}

// ReactionData данные для реакции
type ReactionData struct {
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Emoji      string `json:"emoji"`
}

// SessionStateResponse состояние сессии
type SessionStateResponse struct {
	Session   *entitymodel.Session    `json:"session"`
	Users     []*entitymodel.User     `json:"users"`
	Votes     []*entitymodel.Vote     `json:"votes"`
	Reactions []*entitymodel.Reaction `json:"reactions"`
}
