package websocket

import (
	"backend_go/internal/service"
	"go.uber.org/zap"
)

// BaseHandler содержит общие зависимости для обработчиков
type BaseHandler struct {
	manager *ConnectionManager
	//userService     service.UserService
	sessionService service.SessionService
	//voteService     service.VoteService
	//reactionService service.ReactionService
	log *zap.Logger
}

// NewBaseHandler создает базовый обработчик
func NewBaseHandler(
	manager *ConnectionManager,
	//userService service.UserService,
	sessionService service.SessionService,
	// voteService service.VoteService,
	// reactionService service.ReactionService,
	log *zap.Logger,
) *BaseHandler {
	return &BaseHandler{
		manager: manager,
		//userService:     userService,
		sessionService: sessionService,
		//voteService:     voteService,
		//reactionService: reactionService,
		log: log,
	}
}
