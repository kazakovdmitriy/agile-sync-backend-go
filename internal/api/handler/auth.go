package handler

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/converter"
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
	log         *zap.Logger
}

func NewAuthHandler(authService service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         logger,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req apimodel.UserLogin

	if err := c.ShouldBind(&req); err != nil {
		h.log.Info("Bind Error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.log.Info("Login Error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) GuestLogin(c *gin.Context) {
	h.log.Info("Handle Guest Logout")
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req apimodel.UserRegister

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Info("Bind Error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Debug("Handle Register", zap.String("name", req.Name), zap.String("email", req.Email))

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.log.Info("Register Error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*entitymodel.User) // замените *model.User на ваш реальный тип
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	userProfile := converter.ToUserProfile(user)

	c.JSON(http.StatusOK, userProfile)
}
