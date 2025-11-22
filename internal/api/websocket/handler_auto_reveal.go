package websocket

import (
	"backend_go/internal/model/websocketmodel"
	"backend_go/internal/utils"
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type AutoRevealHandler struct {
	*BaseHandler
}

func NewAutoRevealHandler(baseHandler *BaseHandler) *AutoRevealHandler {
	return &AutoRevealHandler{
		BaseHandler: baseHandler,
	}
}

func (h *AutoRevealHandler) CanHandle(event websocketmodel.SocketEvent) bool {
	return event == websocketmodel.EventAutoReveal
}

func (h *AutoRevealHandler) Handle(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error {
	h.log.Debug("handle auto reveal toggle", zap.Any("data", data))

	var payload websocketmodel.RevealCardsData
	if err := utils.MapToStruct(data, &payload); err != nil {
		h.log.Warn("Invalid reveal cards payload", zap.Any("data", data), zap.Error(err))
		return err
	}

	if err := h.sessionService.AutoRevealCardsInSession(ctx, payload.SessionID); err != nil {
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
