package websocket

import (
	"backend_go/internal/model/websocketmodel"
	"backend_go/internal/utils"
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type RevealCardsHandler struct {
	*BaseHandler
}

func NewRevealCardsHandler(baseHandler *BaseHandler) *RevealCardsHandler {
	return &RevealCardsHandler{
		BaseHandler: baseHandler,
	}
}

func (h *RevealCardsHandler) CanHandle(event websocketmodel.SocketEvent) bool {
	return event == websocketmodel.EventRevealCards
}

func (h *RevealCardsHandler) Handle(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error {
	var payload websocketmodel.RevealCardsData
	if err := utils.MapToStruct(data, &payload); err != nil {
		h.log.Warn("Invalid reveal cards payload", zap.Any("data", data), zap.Error(err))
		return err
	}

	err := h.sessionService.RevealCardsInSession(ctx, payload.SessionID, true)
	if err != nil {
		h.log.Warn("Invalid reveal cards session", zap.Error(err))
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
