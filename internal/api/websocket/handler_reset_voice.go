package websocket

import (
	"backend_go/internal/model/websocketmodel"
	"backend_go/pkg/utils"
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ResetVotesHandler struct {
	*BaseHandler
}

func NewResetVotesHandler(baseHandler *BaseHandler) *ResetVotesHandler {
	return &ResetVotesHandler{
		BaseHandler: baseHandler,
	}
}

func (h *ResetVotesHandler) CanHandle(event websocketmodel.SocketEvent) bool {
	return event == websocketmodel.EventResetVotes
}

func (h *ResetVotesHandler) Handle(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error {
	h.log.Debug("Recive Reset voice handler", zap.Any("data", data))

	var payload websocketmodel.RevealCardsData
	if err := utils.MapToStruct(data, &payload); err != nil {
		h.log.Warn("Invalid reveal cards payload", zap.Any("data", data), zap.Error(err))
		return err
	}

	if err := h.voteService.DeleteVoteInSession(ctx, payload.SessionID); err != nil {
		h.log.Warn("Failed to delete vote in session", zap.Any("session_id", payload.SessionID))
		return err
	}

	if err := h.sessionService.RevealCardsInSession(ctx, payload.SessionID, false); err != nil {
		h.log.Warn("Failed to reveal cards in session", zap.Any("data", payload), zap.Error(err))
		return err
	}

	session, err := h.sessionService.GetSessionByID(ctx, payload.SessionID.String())
	if err != nil {
		h.log.Warn("GetSessionByID", zap.Error(err))
		return err
	}

	broadcast := websocketmodel.BaseMessage{
		Event: string(websocketmodel.EventSessionUpdated),
		Data:  session,
	}
	return h.manager.Broadcast(payload.SessionID.String(), broadcast)
}
