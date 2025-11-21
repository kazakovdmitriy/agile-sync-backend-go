package websocket

import (
	"backend_go/internal/model/websocketmodel"
	"backend_go/internal/utils"
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

func (h *JoinSessionHandler) CanHandle(event websocketmodel.SocketEvent) bool {
	return event == websocketmodel.EventJoinSession
}

func (h *JoinSessionHandler) Handle(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error {
	var payload websocketmodel.JoinSessionData
	if err := utils.MapToStruct(data, &payload); err != nil {
		h.log.Warn("Invalid join_session payload", zap.Any("data", data), zap.Error(err))
		return err
	}

	h.log.Debug("Handle join session",
		zap.String("session_id", payload.SessionID.String()),
		zap.String("user_id", payload.UserID.String()),
		zap.String("user_name", payload.UserName),
		zap.Bool("is_watcher", *payload.IsWatcher),
	)

	err := h.sessionService.ConnectUserToSession(ctx, payload.UserID.String(), payload.SessionID.String())
	if err != nil {
		return err
	}
	h.manager.Connect(payload.SessionID.String(), conn)

	// Отправляем пользователю
	response := websocketmodel.BaseMessage{
		Event: string(websocketmodel.EventJoinSession),
		Data: map[string]interface{}{
			"id":         payload.UserID,
			"name":       payload.UserName,
			"session_id": payload.SessionID,
		},
	}
	if err = h.manager.SendTo(conn, response); err != nil {
		h.log.Error("failed to send join session to user", zap.Any("data", data), zap.Error(err))
		return err
	}

	// Broadcast
	session, err := h.sessionService.GetSessionByID(ctx, payload.SessionID.String())
	if err != nil {
		return err
	}
	broadcast := websocketmodel.BaseMessage{
		Event: string(websocketmodel.EventSessionUpdated),
		Data:  session,
	}
	return h.manager.Broadcast(payload.SessionID.String(), broadcast)
}
