package server

import (
	_ "backend_go/docs"
	"backend_go/internal/api"
	"backend_go/internal/api/handler"
	"backend_go/internal/api/middleware"
	"backend_go/internal/api/websocket"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(
	wsManager *websocket.WebSocketHandler,
	authHandler *handler.AuthHandler,
	sessionHandler *handler.SessionHandler,
	authService api.AuthService,
) *gin.Engine {
	router := gin.Default()

	if gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	apiGroup := router.Group("/api")
	{
		// Auth
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/guest_login", authHandler.GuestLogin)
		}

		authProtectedGroup := apiGroup.Group("/auth")
		authProtectedGroup.Use(middleware.AuthMiddleware(authService))
		{
			authProtectedGroup.GET("/me", authHandler.Me)
		}

		// Session
		sessionGroup := apiGroup.Group("/sessions")
		{
			sessionGroup.GET("/:session_id", sessionHandler.GetSession)
		}
		sessionGroupProtect := apiGroup.Group("/sessions")
		sessionGroupProtect.Use(middleware.AuthMiddleware(authService))
		{
			sessionGroupProtect.GET("", sessionHandler.GetUserSession)
			sessionGroupProtect.POST("", sessionHandler.Create)
			sessionGroupProtect.DELETE("/:session_id", sessionHandler.DeleteSession)
		}
	}

	router.GET("/ws", wsManager.HandleWebSocket)

	return router
}
