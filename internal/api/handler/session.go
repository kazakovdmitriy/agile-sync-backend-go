package handler

import (
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

func (h *SessionHandler) GetUserSession(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*entitymodel.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
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
