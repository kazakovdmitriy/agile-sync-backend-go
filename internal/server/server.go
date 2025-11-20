package server

import (
	"backend_go/internal/api/handler"
	"backend_go/internal/api/websocket"
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
	dbconn, err := db.NewPostgresDB(
		cfg.Database.URL,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		return nil, err
	}

	// Инициализация репозиториев
	userDBRepo := repository.NewUserDBRepo(dbconn.DB, log)
	sessionDBRepo := repository.NewSessionDBRepo(dbconn.DB, log)
	voteDBRepo := repository.NewVoteDBRepo(dbconn.DB, log)

	// Инициализация сервисов
	jwtService := service.NewJwtService(cfg, log)
	authService := service.NewAuthService(userDBRepo, jwtService, log)
	sessionService := service.NewSessionService(sessionDBRepo, voteDBRepo, log)
	voteService := service.NewVoteService(voteDBRepo, log)

	// Инициализация вебсокета
	wsManager := websocket.NewWebSocketHandler(cfg, log, sessionService, voteService)

	// Инициализация хендлеров
	authHandler := handler.NewAuthHandler(authService, log)
	sessionHandler := handler.NewSessionHandler(sessionService, log)

	// Настройка роутинга
	router := setupRouter(wsManager, authHandler, sessionHandler, authService)

	httpServer := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		log:        log,
	}, nil
}

func (s *Server) Run() error {
	s.log.Info("server started", zap.String("addr", s.httpServer.Addr))
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
