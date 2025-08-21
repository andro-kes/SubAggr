package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    ServerPort       string
    GinMode          string
    LogLevel         string
    ShutdownTimeout  time.Duration

    PostgresHost     string
    PostgresPort     string
    PostgresUser     string
    PostgresPassword string
    PostgresDB       string

    AutoMigrate      bool
}

func Load() Config {
    cfg := Config{}

    cfg.ServerPort = getEnvDefault("SERVER_PORT", "8000")
    cfg.GinMode = getEnvDefault("GIN_MODE", "release")
    cfg.LogLevel = getEnvDefault("LOG_LEVEL", "info")

    if v := getEnvDefault("SHUTDOWN_TIMEOUT", "5"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 {
            cfg.ShutdownTimeout = time.Duration(n) * time.Second
        } else {
            cfg.ShutdownTimeout = 5 * time.Second
        }
    }

    cfg.PostgresHost = getEnvDefault("POSTGRES_HOST", "postgres")
    cfg.PostgresPort = getEnvDefault("POSTGRES_PORT", "5432")
    cfg.PostgresUser = getEnvDefault("POSTGRES_USER", "postgres")
    cfg.PostgresPassword = getEnvDefault("POSTGRES_PASSWORD", "postgres")
    cfg.PostgresDB = getEnvDefault("POSTGRES_DB", "subaggr")

    cfg.AutoMigrate = getEnvDefault("AUTO_MIGRATE", "true") != "false"

    return cfg
}

func getEnvDefault(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

