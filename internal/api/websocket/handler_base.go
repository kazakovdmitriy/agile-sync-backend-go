package websocket

import (
	"backend_go/internal/service"
	"go.uber.org/zap"
)

// BaseHandler содержит общие зависимости для обработчиков
type BaseHandler struct {
	manager        *ConnectionManager
	sessionService service.SessionService
	voteService    service.VoteService
	log            *zap.Logger
}

// NewBaseHandler создает базовый обработчик
func NewBaseHandler(
	manager *ConnectionManager,
	sessionService service.SessionService,
	voteService service.VoteService,
	log *zap.Logger,
) *BaseHandler {
	return &BaseHandler{
		manager:        manager,
		sessionService: sessionService,
		voteService:    voteService,
		log:            log,
	}
}
