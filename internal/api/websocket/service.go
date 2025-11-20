package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// WebSocketService маршрутизирует события к соответствующим обработчикам
type WebSocketService struct {
	handlers map[string]EventHandler
	log      *zap.Logger
}

// EventHandler интерфейс для обработчиков событий
type EventHandler interface {
	Handle(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error
	CanHandle(event string) bool
}

// NewWebSocketService создает новый сервис
func NewWebSocketService(baseHandler *BaseHandler, log *zap.Logger) *WebSocketService {
	service := &WebSocketService{
		handlers: make(map[string]EventHandler),
		log:      log,
	}

	// Регистрируем обработчики
	service.RegisterHandler(NewJoinSessionHandler(baseHandler))
	service.RegisterHandler(NewVoteHandler(baseHandler))
	service.RegisterHandler(NewRevealCardsHandler(baseHandler))
	//service.RegisterHandler(sockethandlers.NewResetVotesHandler(baseHandler))
	//service.RegisterHandler(sockethandlers.NewReactionHandler(baseHandler))
	// ... другие обработчики

	return service
}

// RegisterHandler регистрирует обработчик
func (s *WebSocketService) RegisterHandler(handler EventHandler) {
	if joinHandler, ok := handler.(*JoinSessionHandler); ok {
		s.handlers["join_session"] = joinHandler
	} else if voteHandler, ok := handler.(*VoteHandler); ok {
		s.handlers["vote"] = voteHandler
	} else if revealHandler, ok := handler.(*RevealCardsHandler); ok {
		s.handlers["reveal_cards"] = revealHandler
	}
	// TODO: добавить другие типы обработчиков
}

// HandleMessage обрабатывает входящее WebSocket сообщение
func (s *WebSocketService) HandleMessage(ctx context.Context, conn *websocket.Conn, data map[string]interface{}) error {
	event, ok := data["event"].(string)
	if !ok {
		return s.sendError(conn, "event field is required")
	}

	s.log.Info("Processing event", zap.String("event", event))

	handler, exists := s.handlers[event]
	if !exists {
		return s.sendError(conn, "unknown event: "+event)
	}

	return handler.Handle(ctx, conn, data)
}

func (s *WebSocketService) sendError(conn *websocket.Conn, errorMsg string) error {
	return conn.WriteJSON(map[string]interface{}{
		"event": "error",
		"error": errorMsg,
	})
}
