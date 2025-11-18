package websocket

import (
	"github.com/gorilla/websocket"
)

type VoteHandler struct {
	*BaseHandler
}

func NewVoteHandler(baseHandler *BaseHandler) *VoteHandler {
	return &VoteHandler{
		BaseHandler: baseHandler,
	}
}

func (h *VoteHandler) CanHandle(event string) bool {
	return event == "vote"
}

func (h *VoteHandler) Handle(conn *websocket.Conn, data map[string]interface{}) error {
	h.log.Info("Handling vote event")
	return nil
}

//func (h *VoteHandler) Handle(conn *websocket.Conn, data map[string]interface{}) error {
//	sessionID, _ := data["session_id"].(string)
//	userID, _ := data["user_id"].(string)
//	value, _ := data["value"].(string)
//
//	if sessionID == "" || userID == "" || value == "" {
//		return errors.New("session_id, user_id and value are required")
//	}
//
//	// Сохраняем голос
//	vote := service.VoteData{
//		SessionID: sessionID,
//		UserID:    userID,
//		Value:     value,
//	}
//
//	if err := h.voteService.SaveVote(vote); err != nil {
//		return err
//	}
//
//	// Уведомляем о голосовании
//	h.manager.Broadcast(sessionID, map[string]interface{}{
//		"event": "vote_submitted",
//		"data": map[string]interface{}{
//			"user_id":   userID,
//			"has_voted": true,
//		},
//	})
//
//	return nil
//}
