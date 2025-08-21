package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	log.Println("Подключение к базе данных")
	DSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	openDBWithRetry(DSN, 10, 3*time.Second)
}

func openDBWithRetry(dsn string, attempts int, delay time.Duration) {
	for i := 1; i <= attempts; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				log.Printf("Ошибка при получении *sql.DB: %v", err)
			} else {
				sqlDB.SetMaxIdleConns(10)
				sqlDB.SetMaxOpenConns(100)
				sqlDB.SetConnMaxLifetime(time.Hour)
			}

			DB = db

			if os.Getenv("AUTO_MIGRATE") != "false" {
				if err := Migrate(); err != nil {
					log.Printf("Миграция завершилась с ошибкой: %v", err)
				}
			}
			return
		}

		log.Printf("Не удалось подключиться к БД (попытка %d/%d): %v", i, attempts, err)
		time.Sleep(delay)
	}
	log.Fatalf("Не удалось подключиться к базе данных после %d попыток", attempts)
}

func Migrate() error {
	return DB.AutoMigrate(models.Subs{})
}

func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", DB)
		c.Next()
	}
}