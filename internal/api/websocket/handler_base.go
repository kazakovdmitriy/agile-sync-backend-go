package websocket

import (
	"backend_go/internal/api"
	"go.uber.org/zap"
)

// BaseHandler содержит общие зависимости для обработчиков
type BaseHandler struct {
	manager        *ConnectionManager
	sessionService api.SessionService
	voteService    api.VoteService
	userService    api.UserService
	log            *zap.Logger
}

// NewBaseHandler создает базовый обработчик
func NewBaseHandler(
	manager *ConnectionManager,
	sessionService api.SessionService,
	voteService api.VoteService,
	userService api.UserService,
	log *zap.Logger,
) *BaseHandler {
	return &BaseHandler{
		manager:        manager,
		sessionService: sessionService,
		voteService:    voteService,
		userService:    userService,
		log:            log,
	}
}
