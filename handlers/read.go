package handlers

import (
	"log"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/utils"
	"github.com/gin-gonic/gin"
)

func ReadNote(c *gin.Context) {
	var id IDField
	if !utils.CheckError(c, c.ShouldBindJSON(&id), "Не удалось связать данные с полем id") {
		return
	}
	log.Printf("Чтение записи с ID: %d\n", id.ID)

	var sub models.Subs
	if DB := database.GetDB(c); DB != nil {
		obj := DB.Where("id = ?").First(&sub)
		if !utils.CheckError(c, obj.Error, "Не удалось найти запись с таким ID") {
			return
		}
	}

	c.JSON(200, sub)
}