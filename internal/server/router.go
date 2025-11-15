package server

import (
	"backend_go/internal/api/handler"
	"backend_go/internal/api/middleware"
	"backend_go/internal/service"
	"github.com/gin-gonic/gin"
)

func setupRouter(
	authHandler *handler.AuthHandler,
	sessionHandler *handler.SessionHandler,
	authService service.AuthService,
) *gin.Engine {
	router := gin.Default()

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
		sessionGroup.Use(middleware.AuthMiddleware(authService))
		{
			sessionGroup.GET("", sessionHandler.GetUserSession)
			sessionGroup.POST("", sessionHandler.Create)
		}
	}

	return router
}
