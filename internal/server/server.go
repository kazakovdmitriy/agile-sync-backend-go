package server

import (
	"backend_go/internal/api/handler"
	"backend_go/internal/api/websocket"
	"backend_go/internal/infrastructure/config"
	"backend_go/internal/infrastructure/db"
	"backend_go/internal/repository"
	"backend_go/internal/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	bgCancel   context.CancelFunc
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
	userService := service.NewUserService(userDBRepo, log)
	sessionService := service.NewSessionService(sessionDBRepo, voteDBRepo, userDBRepo, log)
	voteService := service.NewVoteService(voteDBRepo, log)

	// Инициализация вебсокета
	wsManager := websocket.NewWebSocketHandler(cfg, log, userService, sessionService, voteService)

	// Инициализация хендлеров
	authHandler := handler.NewAuthHandler(authService, cfg, log)
	sessionHandler := handler.NewSessionHandler(sessionService, log)

	// Настройка роутинга
	router := setupRouter(wsManager, authHandler, sessionHandler, authService)

	bgCtx, bgCancel := context.WithCancel(context.Background())

	go startBackgroundJobs(bgCtx, cfg, userDBRepo, log)

	httpServer := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		log:        log,
		bgCancel:   bgCancel,
	}, nil
}

func (s *Server) Run() error {
	s.log.Info("server started", zap.String("addr", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.bgCancel()
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

func startBackgroundJobs(ctx context.Context, cfg *config.Config, userRepo repository.UserRepository, log *zap.Logger) {
	if cfg.GuestCleanInterval <= 0 {
		log.Info("Guest cleanup disabled (interval <= 0)")
		return
	}

	log.Info("Starting guest cleanup job",
		zap.Duration("interval", cfg.GuestCleanInterval))

	// Первый запуск сразу после старта (опционально)
	cleanInactiveGuests(ctx, cfg, userRepo, log)

	// Планируем регулярные запуски
	ticker := time.NewTicker(cfg.GuestCleanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanInactiveGuests(ctx, cfg, userRepo, log)
		case <-ctx.Done():
			log.Info("Stopping background jobs...")
			return
		}
	}
}

func cleanInactiveGuests(ctx context.Context, cfg *config.Config, userRepo repository.UserRepository, log *zap.Logger) {
	// Используем контекст для таймаута операции
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	intervalStr := durationToPostgresInterval(cfg.GuestCleanInterval)

	deleted, err := userRepo.DeleteInactiveGuests(ctx, intervalStr)
	if err != nil {
		log.Error("Failed to clean inactive guests", zap.Error(err))
		return
	}

	if deleted > 0 {
		log.Info("Cleaned inactive guests", zap.Int64("count", deleted))
	} else {
		log.Info("Cleaned inactive guests", zap.Int64("count", 0))
	}
}

func durationToPostgresInterval(d time.Duration) string {
	// Используем total nanoseconds для точности
	return fmt.Sprintf("%d microseconds", d.Microseconds())
}
