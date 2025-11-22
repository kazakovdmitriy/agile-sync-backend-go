package websocket

import (
	"backend_go/internal/api"
	"backend_go/internal/infrastructure/config"
	"go.uber.org/zap"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler обрабатывает HTTP запросы и управляет жизненным циклом соединений
type WebSocketHandler struct {
	manager        *ConnectionManager
	service        *WebSocketService
	sessionService api.SessionService
	config         *config.Config
	log            *zap.Logger
}

// NewWebSocketHandler создает новый WebSocket handler
func NewWebSocketHandler(
	cfg *config.Config,
	log *zap.Logger,
	userService api.UserService,
	sessionService api.SessionService,
	voteService api.VoteService,
) *WebSocketHandler {
	manager := NewConnectionManager(cfg, log)

	// Создаем базовый обработчик с общими зависимостями
	baseHandler := NewBaseHandler(manager, sessionService, voteService, userService, log)

	// Создаем сервис с зарегистрированными обработчиками
	service := NewWebSocketService(baseHandler, log)

	return &WebSocketHandler{
		manager:        manager,
		service:        service,
		sessionService: sessionService,
		config:         cfg,
		log:            log,
	}
}

// HandleWebSocket обработчик для Gin
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	upgrader := h.manager.GetUpgrader()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("upgrader error", zap.Error(err))
		return
	}
	defer conn.Close()

	h.log.Info("new websocket connection", zap.String("From", conn.RemoteAddr().String()))

	h.configureConnection(conn)

	var sessionID string
	done := make(chan struct{})

	// Запускаем ping handler
	go h.startPingHandler(conn, done)

	// Основной цикл обработки сообщений
	for {
		var data map[string]interface{}
		if err := conn.ReadJSON(&data); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.log.Error("WebSocket error", zap.Error(err))
			}
			break
		}

		h.log.Info("WebSocket received", zap.Any("data", data))

		// Обрабатываем сообщение через сервис
		if err := h.service.HandleMessage(c.Request.Context(), conn, data); err != nil {
			h.log.Error("WebSocket error", zap.Error(err))
		}

		// Сохраняем sessionID для disconnect
		if event, _ := data["event"].(string); event == "join_session" {
			if sid, ok := data["session_id"].(string); ok {
				sessionID = sid
			}
		}
	}

	close(done)

	// Отключаем от сессии при закрытии соединения
	if sessionID != "" {
		userID, hasUserID := h.manager.GetUserIDByConnection(conn)
		h.manager.Disconnect(sessionID, conn)

		// Если есть userID, обновляем состояние в БД и отправляем state_update
		if hasUserID && userID != "" {
			h.log.Info("User disconnected, updating session state",
				zap.String("userID", userID),
				zap.String("sessionID", sessionID),
			)

			// Обновляем состояние пользователя в БД
			if err := h.sessionService.DisconnectUser(c.Request.Context(), userID, sessionID); err != nil {
				h.log.Warn("Failed to disconnect user from session in DB",
					zap.String("userID", userID),
					zap.String("sessionID", sessionID),
					zap.Error(err),
				)
			} else {
				// Отправляем обновленное состояние сессии остальным пользователям
				session, err := h.sessionService.GetSessionByID(c.Request.Context(), sessionID)
				if err != nil {
					h.log.Warn("Failed to get session for state update",
						zap.String("sessionID", sessionID),
						zap.Error(err),
					)
				} else {
					broadcast := map[string]interface{}{
						"event": "state_update",
						"data":  session,
					}
					if err := h.manager.Broadcast(sessionID, broadcast); err != nil {
						h.log.Warn("Failed to broadcast state update",
							zap.String("sessionID", sessionID),
							zap.Error(err),
						)
					}
				}
			}
		} else {
			h.log.Warn("User disconnected but userID not found in connection mapping",
				zap.String("sessionID", sessionID),
			)
		}
	}
}

// configureConnection настраивает WebSocket соединение
func (h *WebSocketHandler) configureConnection(conn *websocket.Conn) {
	conn.SetReadLimit(h.config.WebSocket.MaxMessageSize)

	if h.config.WebSocket.PongTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(h.config.WebSocket.PongTimeout))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(h.config.WebSocket.PongTimeout))
			return nil
		})
	}
}

// startPingHandler запускает горутину для отправки ping сообщений
func (h *WebSocketHandler) startPingHandler(conn *websocket.Conn, done <-chan struct{}) {
	if h.config.WebSocket.PingInterval <= 0 {
		return
	}

	ticker := time.NewTicker(h.config.WebSocket.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(h.config.WebSocket.WriteTimeout)); err != nil {
				h.log.Error("WebSocket ping error", zap.Error(err))
				return
			}
		case <-done:
			return
		}
	}
}

// GetManager возвращает менеджер для использования в других частях приложения
func (h *WebSocketHandler) GetManager() *ConnectionManager {
	return h.manager
}

// HealthHandler возвращает метрики WebSocket
func (h *WebSocketHandler) HealthHandler(c *gin.Context) {
	manager := h.GetManager()

	c.JSON(200, gin.H{
		"status":            "healthy",
		"active_sessions":   manager.GetSessionCount(),
		"total_connections": manager.GetConnectionCount(),
	})
}
