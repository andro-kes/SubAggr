package app

import (
	"log"
	"os"

	"github.com/andro-kes/SubAggr/internal/handlers"
	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/gin-gonic/gin"
)

func Run() {
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	router := gin.Default()
	router.Use(database.DBMiddleware())

	router.POST("/SUBS", handlers.CreateNote)
	router.DELETE("/SUBS/:id", handlers.DeleteNote)
	router.PUT("/SUBS/:id", handlers.UpdateNote)
	router.GET("/SUBS/:id", handlers.ReadNote)
	router.GET("/SUBS", handlers.ListNotes)
	router.POST("/SUBS/SUMMARY", handlers.SumPriceSubs)

	registerSwagger(router)

	if err := router.Run(":8000"); err != nil {
		log.Fatalln("Не удалось запустить сервер")
	}
}