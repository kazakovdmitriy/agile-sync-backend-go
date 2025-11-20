package websocket

import (
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/model/websocketmodel"
	"backend_go/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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

func (h *VoteHandler) Handle(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error {
	var payload websocketmodel.VoteData
	if err := utils.MapToStruct(data, &payload); err != nil {
		h.log.Warn("Invalid vote payload", zap.Any("data", data), zap.Error(err))
		return fmt.Errorf("invalid vote data: %w", err)
	}

	// Валидация UUID (опционально — можно вынести в отдельную функцию)
	sessionUUID, err := uuid.Parse(payload.SessionID.String())
	if err != nil {
		h.log.Warn("Invalid session_id format", zap.String("session_id", payload.SessionID.String()), zap.Error(err))
		return errors.New("invalid session_id format")
	}

	userUUID, err := uuid.Parse(payload.UserID.String())
	if err != nil {
		h.log.Warn("Invalid user_id format", zap.String("user_id", payload.UserID.String()), zap.Error(err))
		return errors.New("invalid user_id format")
	}

	vote := entitymodel.Vote{
		SessionID: sessionUUID,
		UserID:    userUUID,
		Value:     payload.Value,
	}

	voteUUID, err := h.voteService.SaveVote(ctx, &vote)
	if err != nil {
		h.log.Error("Failed to save vote", zap.Error(err),
			zap.String("session_id", payload.SessionID.String()),
			zap.String("user_id", payload.UserID.String()),
			zap.String("vote_id", voteUUID.String()),
		)
		return fmt.Errorf("failed to save vote: %w", err)
	}

	// Получаем обновлённую сессию (если нужно для broadcast)
	session, err := h.sessionService.GetSessionByID(ctx, payload.SessionID.String())
	if err != nil {
		h.log.Error("Failed to get session after vote", zap.Error(err),
			zap.String("session_id", payload.SessionID.String()))
	}

	// Отправляем подтверждение пользователю
	response := websocketmodel.BaseMessage{
		Event: string(websocketmodel.EventVote),
		Data: map[string]interface{}{
			"id":         voteUUID.String(),
			"session_id": payload.SessionID,
			"user_id":    payload.UserID,
			"value":      payload.Value,
		},
	}
	if err := h.manager.SendTo(conn, response); err != nil {
		h.log.Error("Failed to send vote confirmation", zap.Error(err))
	}

	// Broadcast обновлённого состояния сессии (если session != nil)
	if session != nil {
		broadcast := websocketmodel.BaseMessage{
			Event: string(websocketmodel.EventSessionUpdated),
			Data:  session,
		}
		if err := h.manager.Broadcast(payload.SessionID.String(), broadcast); err != nil {
			h.log.Error("Failed to broadcast session update", zap.Error(err))
		}
	}

	return nil
}
