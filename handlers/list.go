package handlers

import (
	"log"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/gin-gonic/gin"
)

func ListNotes(c *gin.Context) {
	var notes []models.Subs
	if DB := database.GetDB(c); DB != nil {
		DB.Find(&notes)
	}
	if len(notes) == 0 {
		log.Println("Ни одной задачи не было найдено")
		c.JSON(200, gin.H{"Status": "Не найдено"})
		return
	}
	c.JSON(200, notes)
}