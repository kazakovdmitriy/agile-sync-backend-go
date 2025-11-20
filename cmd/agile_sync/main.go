package main

import (
	"backend_go/internal/infrastructure/config"
	"backend_go/internal/infrastructure/logger"
	server2 "backend_go/internal/server"
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/// main.go
// @title Backend Go API
// @version 1.0
// @description API для управления сессиями планирования с системой аутентификации и WebSocket
// @BasePath /api

// @contact.name Dmitriy Kazakov
// @contact.email kda94@mail.ru

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT токен в формате: "Bearer {token}"

// @tag.name auth
// @tag.description Аутентификация и управление пользователями

// @tag.name sessions
// @tag.description Управление сессиями планирования и голосования
func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.LoadConfig("configs/config.yaml")

	fmt.Println(cfg)

	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Starting server...")

	server, err := server2.NewServer(cfg, logger)
	if err != nil {
		logger.Fatal("failed to initialize server", zap.Error(err))
	}

	go func() {
		if err := server.Run(); err != nil {
			logger.Fatal("failed to run server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Fatal("failed to shutting down server", zap.Error(err))
	}
}
