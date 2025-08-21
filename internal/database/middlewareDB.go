package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(host, port, user, password, dbname string, autoMigrate bool) error {
	log.Println("Подключение к базе данных")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbname,
		port,
	)

	if err := openDBWithRetry(dsn, 10, 3*time.Second); err != nil {
		return err
	}

	if autoMigrate {
		if err := Migrate(); err != nil {
			log.Printf("Миграция завершилась с ошибкой: %v", err)
		}
	}
	return nil
}

func openDBWithRetry(dsn string, attempts int, delay time.Duration) error {
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
			return nil
		}

		log.Printf("Не удалось подключиться к БД (попытка %d/%d): %v", i, attempts, err)
		time.Sleep(delay)
	}
	return fmt.Errorf("не удалось подключиться к базе данных после %d попыток", attempts)
}

func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", DB)
		c.Next()
	}
}