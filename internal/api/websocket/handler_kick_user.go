package websocket

import (
	"backend_go/internal/model/websocketmodel"
	"backend_go/internal/utils"
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type KickUserHandler struct {
	*BaseHandler
}

func NewKickUserHandler(baseHandler *BaseHandler) *KickUserHandler {
	return &KickUserHandler{
		BaseHandler: baseHandler,
	}
}

func (h *KickUserHandler) CanHandle(event websocketmodel.SocketEvent) bool {
	return event == websocketmodel.EventKickUser
}

func (h *KickUserHandler) Handle(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error {
	var payload websocketmodel.KickUserData
	if err := utils.MapToStruct(data, &payload); err != nil {
		h.log.Warn("Invalid kick_user payload", zap.Any("data", data), zap.Error(err))
		return err
	}

	h.log.Debug("Handle kick user",
		zap.String("session_id", payload.SessionID.String()),
		zap.String("initiator_id", payload.InitiatorUserID.String()),
		zap.String("target_user", payload.TargetUser.String()),
	)

	session, err := h.sessionService.GetSessionByID(ctx, payload.SessionID.String())
	if err != nil {
		h.log.Warn("Failed to get session", zap.Error(err))
		return err
	}

	if session.CreatorID != *payload.InitiatorUserID {
		response := map[string]interface{}{
			"event":   websocketmodel.EventKickUserError,
			"message": "Только создатель сессии может удалять пользователей",
		}

		if err := h.manager.SendTo(conn, response); err != nil {
			h.log.Warn("Failed to send message", zap.Error(err))
			return err
		}

		return nil
	}

	kickedUser, err := h.userService.GetUser(ctx, *payload.TargetUser)
	if err != nil {
		h.log.Warn("Failed to get kicked user", zap.Error(err))
		return err
	}

	err = h.sessionService.DisconnectUser(ctx, kickedUser.ID.String(), payload.SessionID.String())
	if err != nil {
		h.log.Warn("Failed to disconnect kicked user", zap.Error(err))
		return err
	}

	kickedResponse := map[string]interface{}{
		"event":   websocketmodel.EventKicked,
		"message": "Вы были удалены из сессии",
	}
	if err := h.manager.SendToUser(payload.TargetUser.String(), kickedResponse); err != nil {
		h.log.Warn("Failed to send kicked event to user", zap.Error(err))
	}

	response := map[string]interface{}{
		"event":            websocketmodel.EventUserKickedBroadcast,
		"kicked_user_id":   payload.TargetUser.String(),
		"kicked_user_name": kickedUser.Name,
	}
	if err := h.manager.Broadcast(payload.SessionID.String(), response); err != nil {
		h.log.Warn("Failed to send message", zap.Error(err))
		return err
	}

	h.manager.DisconnectUser(payload.TargetUser.String())

	session, err = h.sessionService.GetSessionByID(ctx, payload.SessionID.String())
	if err != nil {
		return err
	}
	broadcast := websocketmodel.BaseMessage{
		Event: string(websocketmodel.EventSessionUpdated),
		Data:  session,
	}
	err = h.manager.Broadcast(payload.SessionID.String(), broadcast)
	if err != nil {
		h.log.Warn("Failed to send message", zap.Error(err))
		return err
	}

	response = map[string]interface{}{
		"event":            websocketmodel.EventUserKicked,
		"kicked_user_name": kickedUser.Name,
	}
	return h.manager.SendTo(conn, response)
}
