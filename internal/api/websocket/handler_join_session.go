package websocket

import (
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

func (h *JoinSessionHandler) CanHandle(event string) bool {
	return event == "join_session"
}

func (h *JoinSessionHandler) Handle(conn *websocket.Conn, data map[string]interface{}) error {
	h.log.Debug("Handle join session", zap.Any("data", data))

	err := h.manager.SendTo(conn, map[string]interface{}{
		"event": "join_session",
		"user": map[string]interface{}{
			"id":         data["user_id"],
			"name":       data["name"],
			"session_id": data["session_id"],
		},
	})
	if err != nil {
		h.log.Error("failed to send join session to user", zap.Any("data", data), zap.Error(err))
		return err
	}

	return nil
}

//func (h *JoinSessionHandler) Handle(conn *websocket.Conn, data map[string]interface{}) error {
//	sessionID, _ := data["session_id"].(string)
//	userID, _ := data["user_id"].(string)
//	userName, _ := data["user_name"].(string)
//
//	if sessionID == "" || userID == "" {
//		return errors.New("session_id and user_id are required")
//	}
//
//	// Подключаем к сессии
//	h.manager.Connect(sessionID, conn)
//
//	// Бизнес-логика: создаем/получаем пользователя
//	//user, err := h.userService.GetOrCreateUser(userID, userName, sessionID)
//	//if err != nil {
//	//	return err
//	//}
//
//	// Получаем состояние сессии
//	sessionState, err := h.getSessionState(sessionID)
//	if err != nil {
//		return err
//	}
//
//	// Отправляем состояние сессии клиенту
//	h.manager.SendTo(conn, map[string]interface{}{
//		"event": "session_state",
//		"data":  sessionState,
//	})
//
//	// Уведомляем других участников
//	h.manager.Broadcast(sessionID, map[string]interface{}{
//		"event": "user_joined",
//		"data": map[string]interface{}{
//			"user_id":    userID,
//			"user_name":  userName,
//			"is_watcher": user.IsWatcher,
//		},
//	})
//
//	return nil
//}

//func (h *JoinSessionHandler) getSessionState(sessionID string) (map[string]interface{}, error) {
//	// Здесь получаем состояние сессии из сервисов
//	// Это примерная реализация
//	session, err := h.sessionService.GetByID(sessionID)
//	if err != nil {
//		return nil, err
//	}
//
//	users, err := h.userService.GetBySessionID(sessionID)
//	if err != nil {
//		return nil, err
//	}
//
//	return map[string]interface{}{
//		"session": session,
//		"users":   users,
//	}, nil
//}
