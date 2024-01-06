// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"nox-ai/internal/delivery"
	"nox-ai/internal/repository"
	"nox-ai/internal/service"
	"nox-ai/internal/usecase"
	"nox-ai/pkg/client"
	"nox-ai/pkg/config"
	"nox-ai/pkg/migration"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// Load Environment
	config.Environment()

	// Initial configuration
	cfg := config.InitService()

	// Initial logger
	logger := cfg.GetLogger()

	// Init Database
	logger.Info("Initializing database connection")
	db := cfg.GetDatabase()

	// Init Redis
	logger.Info("Initializing redis connection")
	redis := cfg.GetRedis()

	// Init GPT
	logger.Info("Initializing gpt client")
	gpt := cfg.GetGPT()

	// Automigrate Database
	migration.AutoMigrate(db)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	ctx, cancel := context.WithCancel(context.Background())

	initRoute(e, db, redis, gpt, logger, cfg.GetWhatsapp())

	serverErr := make(chan os.Signal, 1)
	signal.Notify(serverErr, os.Interrupt)

	go func() {
		e.Logger.Info("Server started")
		if err := e.Start(fmt.Sprintf(":%v", cfg.GetServicePort())); err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down server")
		}
	}()

	select {
	case <-serverErr:
		e.Logger.Print("Shutting down server gracefully...")

		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		if err := e.Shutdown(shutdownCtx); err != nil {
			e.Logger.Printf("Server shutdown error: %v", err)
		}

		e.Logger.Info("Server gracefully stopped")
		cancel()
	case <-ctx.Done():
		e.Logger.Info("Server stopped")
	}

}

func initRoute(e *echo.Echo, db *gorm.DB, redis *redis.Client, gpt *openai.Client, logger *zap.Logger, wa config.Whatsapp) {

	repo := repository.NewRepository(db)
	service := service.NewService(gpt, client.NewHttpClient("https://graph.facebook.com", context.Background(), logger, 10*time.Second), wa)
	usecase := usecase.NewUsecase(repo, redis, logger, service)
	handler := delivery.NewDelivery(e, usecase, logger)

	e.POST("/message", handler.Message)
}
