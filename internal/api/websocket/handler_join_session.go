package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type JoinSessionHandler struct {
	*BaseHandler
}

func NewJoinSessionHandler(baseHandler *BaseHandler) *JoinSessionHandler {
	return &JoinSessionHandler{
		BaseHandler: baseHandler,
	}
}

func (h *JoinSessionHandler) CanHandle(event string) bool {
	return event == "join_session"
}

func (h *JoinSessionHandler) Handle(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error {
	h.log.Debug("Handle join session", zap.Any("data", data))

	// TODO: Мне не нравится преобразование мапы в строку, нужно сделать типобезопасную конвертацию
	sessionID := data["session_id"].(string)
	userID := data["user_id"].(string)
	userName := data["user_name"].(string)
	//isWatcher := data["is_watcher"].(bool)

	err := h.sessionService.ConnectUserToSession(ctx, userID, sessionID)
	if err != nil {
		return err
	}
	h.manager.Connect(sessionID, conn)

	// Отправляем пользователю
	err = h.manager.SendTo(conn, map[string]interface{}{
		"event": "join_session",
		"user": map[string]interface{}{
			"id":         userID,
			"name":       userName,
			"session_id": sessionID,
		},
	})
	if err != nil {
		h.log.Error("failed to send join session to user", zap.Any("data", data), zap.Error(err))
		return err
	}

	session, err := h.sessionService.GetSessionByID(ctx, sessionID)
	if err != nil {
		h.log.Error("failed to get session by id", zap.Any("data", data), zap.Error(err))
		return err
	}

	err = h.manager.Broadcast(sessionID, map[string]interface{}{
		"event":   "state_update",
		"session": session,
	})
	if err != nil {
		h.log.Error("failed to broadcast session to user", zap.Any("data", data), zap.Error(err))
		return err
	}

	return nil
}
