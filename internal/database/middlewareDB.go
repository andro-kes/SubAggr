package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/utils"
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
	
	openDB(DSN)
}

func openDB(dsn string) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	utils.MustNotError(err, "Не удалось открыть базу данных")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка при получении *sql.DB: %v", err)
	}

    sqlDB.SetMaxIdleConns(10)  
    sqlDB.SetMaxOpenConns(100)   
    sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	utils.MustNotError(Migrate(), "Не удалось выполнить миграции")
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