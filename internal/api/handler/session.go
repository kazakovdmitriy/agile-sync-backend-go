package handler

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type SessionHandler struct {
	sessionService service.SessionService
	log            *zap.Logger
}

func NewSessionHandler(sessionService service.SessionService, log *zap.Logger) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		log:            log,
	}
}

// Create создание новой сессии голосования
// @Summary Создать сессию
// @Description Создание новой сессии для планирования с возможностью голосования
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body apimodel.SessionCreate true "Данные для создания сессии"
// @Success 201 {object} apimodel.CreateSessionResponse "Сессия успешно создана"
// @Failure 400 {object} ErrorResponse "Ошибка валидации данных"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /sessions [post]
func (h *SessionHandler) Create(c *gin.Context) {
	var req apimodel.SessionCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Info("Bind Error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := h.getUser(c)
	if !ok {
		return
	}

	session, err := h.sessionService.CreateSession(c.Request.Context(), &req, user)
	if err != nil {
		h.log.Info("Create Error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"session_id": session.ID.String(),
			"user_id":    user.ID.String(),
			"created_at": session.CreatedAt,
		},
	)
}

// GetUserSession получение сессий текущего пользователя
// @Summary Получить мои сессии
// @Description Получение списка сессий, созданных текущим пользователем
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entitymodel.Session "Список сессий пользователя"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 403 {object} ErrorResponse "Гостевой пользователь не имеет доступа"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /sessions [get]
func (h *SessionHandler) GetUserSession(c *gin.Context) {
	user, ok := h.getUser(c)
	if !ok {
		return
	}

	if user.IsGuest {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"error": "Гостевые пользователи не могут просматривать свои сессии."},
		)
		return
	}

	sessions, err := h.sessionService.GetUserSession(c.Request.Context(), user.ID.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting user sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// GetSession получение сессии по ID
// @Summary Получить сессию по ID
// @Description Получение информации о конкретной сессии по её идентификатору
// @Tags sessions
// @Accept json
// @Produce json
// @Param session_id path string true "UUID сессии"
// @Success 200 {object} entitymodel.Session "Данные сессии"
// @Failure 404 {object} ErrorResponse "Сессия не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /sessions/{session_id} [get]
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	session, err := h.sessionService.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error getting session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// DeleteSession удаление сессии
// @Summary Удалить сессию
// @Description Удаление сессии по ID (только для создателя сессии)
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path string true "UUID сессии"
// @Success 200 {object} apimodel.DeleteSessionResponse "Сессия успешно удалена"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 403 {object} ErrorResponse "Доступ запрещен (гостевой пользователь или не создатель сессии)"
// @Failure 404 {object} ErrorResponse "Сессия не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /sessions/{session_id} [delete]
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	user, ok := h.getUser(c)
	if !ok {
		return
	}

	if user.IsGuest {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"detail": "Гостевые пользователи не могут удалять сессии."},
		)
		return
	}

	sessionID := c.Param("session_id")
	err := h.sessionService.DeleteSession(c.Request.Context(), sessionID, user.ID.String())
	if err != nil {
		h.log.Info("Delete Error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error deleting session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Сессия успешно удалена"})
}

func (h *SessionHandler) getUser(c *gin.Context) (*entitymodel.User, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return nil, false
	}

	user, ok := userInterface.(*entitymodel.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return nil, false
	}

	return user, true
}
