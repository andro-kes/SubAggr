package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andro-kes/SubAggr/internal/config"
	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Run() {
	cfg := config.Load()
	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	if err := database.Init(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.AutoMigrate); err != nil {
		slog.Error("DB init failed", slog.String("error", err.Error()))
		return
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(database.DBMiddleware())

	router.POST("/SUBS", handlers.CreateNote)
	router.DELETE("/SUBS/:id", handlers.DeleteNote)
	router.PUT("/SUBS/:id", handlers.UpdateNote)
	router.GET("/SUBS/:id", handlers.ReadNote)
	router.GET("/SUBS", handlers.ListNotes)
	router.POST("/SUBS/SUMMARY", handlers.SumPriceSubs)

	registerSwagger(router)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen error", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
	}
}