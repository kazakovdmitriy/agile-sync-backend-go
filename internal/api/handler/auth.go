package handler

import (
	"backend_go/internal/model/apimodel"
	"backend_go/internal/model/converter"
	"backend_go/internal/model/entitymodel"
	"backend_go/internal/service"
	"errors"
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

// Login аутентификация пользователя
// @Summary Вход в систему
// @Description Аутентификация пользователя по email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param input body apimodel.UserLogin true "Данные для входа"
// @Success 200 {object} apimodel.TokenResponse "Успешная аутентификация"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 401 {object} ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/login [post]
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

// GuestLogin гостевой логин
// @Summary Позволяет войти как гость
// @Description Создание временной гостевой учетной записи
// @Tags auth
// @Accept json
// @Produce json
// @Param input body apimodel.GuestLogin true "Данные гостя"
// @Success 200 {object} apimodel.TokenResponse "Успешная аутентификация гостя"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/guest_login [post]
func (h *AuthHandler) GuestLogin(c *gin.Context) {
	var req apimodel.GuestLogin
	if err := c.ShouldBind(&req); err != nil {
		h.log.Info("Bind Error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error from server"})
		return
	}

	resp, err := h.authService.GuestLogin(c.Request.Context(), &req)
	if err != nil {
		h.log.Info("GuestLogin Error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GuestLogin failure"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Register регистрация пользователя
// @Summary Регистрация пользователя
// @Description Создание нового аккаунта пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body apimodel.UserRegister true "Данные для регистрации"
// @Success 200 {object} apimodel.TokenResponse "Успешная регистрация"
// @Failure 400 {object} ErrorResponse "Пользователь уже существует или ошибка валидации"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req apimodel.UserRegister

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Info("Bind Error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.log.Info("Register Error", zap.Error(err))

		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "Пользователь с таким email уже существует"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Me получение информации о текущем пользователе
// @Summary Получить текущего пользователя
// @Description Получение информации об авторизованном пользователе
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} apimodel.UserProfile "Данные пользователя"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
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

	userProfile := converter.ToUserProfile(user)

	c.JSON(http.StatusOK, userProfile)
}
