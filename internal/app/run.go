package app

import (
	"log"

	_ "github.com/andro-kes/SubAggr/docs" // импорт сгенерированной документации
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/andro-kes/SubAggr/internal/handlers"
	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()
	router.Use(database.DBMiddleware())

	router.POST("/SUBS", handlers.CreateNote)
	router.DELETE("/SUBS/:id", handlers.DeleteNote)
	router.PUT("/SUBS/:id", handlers.UpdateNote)
	router.GET("/SUBS/:id", handlers.ReadNote)
	router.GET("/SUBS", handlers.ListNotes)
	router.POST("/SUBS/SUMMARY", handlers.SumPriceSubs)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8000"); err != nil {
		log.Fatalln("Не удалось запустить сервер")
	}
}