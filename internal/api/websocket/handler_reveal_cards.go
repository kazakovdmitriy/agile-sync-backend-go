package websocket

import (
	"context"
	"github.com/gorilla/websocket"
)

type RevealCardsHandler struct {
	*BaseHandler
}

func NewRevealCardsHandler(baseHandler *BaseHandler) *RevealCardsHandler {
	return &RevealCardsHandler{
		BaseHandler: baseHandler,
	}
}

func (h *RevealCardsHandler) CanHandle(event string) bool {
	return event == "reveal_cards"
}

func (h *RevealCardsHandler) Handle(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error {
	h.log.Info("Handling reveal card event")
	return nil
}

//func (h *RevealCardsHandler) Handle(conn *websocket.Conn, data map[string]interface{}) error {
//	sessionID, _ := data["session_id"].(string)
//	userID, _ := data["user_id"].(string)
//
//	if sessionID == "" {
//		return errors.New("session_id is required")
//	}
//
//	// Проверяем права
//	session, err := h.sessionService.GetByID(sessionID)
//	if err != nil {
//		return err
//	}
//
//	if session.CreatorID != userID {
//		return errors.New("only session creator can reveal cards")
//	}
//
//	// Обновляем сессию
//	if err := h.sessionService.RevealCards(sessionID); err != nil {
//		return err
//	}
//
//	// Получаем голоса и рассылаем
//	votes, err := h.voteService.GetBySessionID(sessionID)
//	if err != nil {
//		return err
//	}
//
//	h.manager.Broadcast(sessionID, map[string]interface{}{
//		"event": "cards_revealed",
//		"data": map[string]interface{}{
//			"votes": votes,
//		},
//	})
//
//	return nil
//}
