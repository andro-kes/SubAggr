package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/andro-kes/SubAggr/internal/app"
)

// @title Subscriptions API
// @version 1.0
// @description API для управления подписками пользователей в Ed-tech стартапе (агрегация и аналитика подписок)
// @host localhost:8000
// @BasePath /
func main() {
	setupLogger()
	app.Run()
}

func setupLogger() {
	level := os.Getenv("LOG_LEVEL")
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	slog.SetDefault(slog.New(handler))
}