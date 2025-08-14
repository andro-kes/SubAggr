package handlers

import (
	"log"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/utils"
	"github.com/gin-gonic/gin"
)

func UpdateNote(c *gin.Context) {
	var updatedSub models.Subs
	if !utils.CheckError(c, c.ShouldBindJSON(&updatedSub), "Невалидные данные для связывания записи") {
		return
	}
	log.Printf("Обновление записи с ID: %d\n", updatedSub.ID)

	if DB := database.GetDB(c); DB != nil {
		obj := DB.Save(&updatedSub)
		if !utils.CheckError(c, obj.Error, "Не удалось обновить данные записи") {
			return
		}
	}

	c.JSON(200, updatedSub)
}