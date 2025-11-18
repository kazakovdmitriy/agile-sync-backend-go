package websocketmodel

import "backend_go/internal/model/entitymodel"

// BaseMessage базовое сообщение WebSocket
type BaseMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// JoinSessionData данные для присоединения к сессии
type JoinSessionData struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	IsWatcher bool   `json:"is_watcher"`
}

// VoteData данные для голосования
type VoteData struct {
	UserID string `json:"user_id"`
	Value  string `json:"value"`
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
