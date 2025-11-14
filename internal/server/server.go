package server

import (
	"backend_go/internal/api/handler"
	"backend_go/internal/api/middleware"
	"backend_go/internal/infrastructure/config"
	"backend_go/internal/infrastructure/db"
	"backend_go/internal/repository"
	"backend_go/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	log        *zap.Logger
}

func NewServer(cfg *config.Config, log *zap.Logger) (*Server, error) {

	setupGin(cfg)

	// Инициализация БД
	dbconn, err := db.NewPostgresDB(cfg.DatabaseURL, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	if err != nil {
		return nil, err
	}

	// Инициализация репозиториев
	userDBRepo := repository.NewUserDBRepo(dbconn.DB, log)
	sessionDBRepo := repository.NewSessionDBRepo(dbconn.DB, log)

	// Инициализация сервисов
	jwtService := service.NewJwtService(cfg, log)
	authService := service.NewAuthService(userDBRepo, jwtService, log)
	sessionService := service.NewSessionService(sessionDBRepo, log)

	// Инициализация хендлеров
	authHandler := handler.NewAuthHandler(authService, log)
	sessionHandler := handler.NewSessionHandler(sessionService, log)

	// Настройка роутинга
	router := setupRouter(authHandler, sessionHandler, authService)

	httpServer := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		log:        log,
	}, nil
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func setupGin(cfg *config.Config) {
	if cfg.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	}

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
}

func setupRouter(
	authHandler *handler.AuthHandler,
	sessionHandler *handler.SessionHandler,
	authService service.AuthService,
) *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api")
	{
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

		sessionGroup := apiGroup.Group("/sessions")
		sessionGroup.Use(middleware.AuthMiddleware(authService))
		{
			sessionGroup.GET("", sessionHandler.GetUserSession)
		}
	}

	return router
}
