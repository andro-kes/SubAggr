package handlers

import (
	"log"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/utils"
	"github.com/gin-gonic/gin"
)

type IDField struct {
	ID uint
}

func DeleteNote(c *gin.Context) {
	var id IDField
	if !utils.CheckError(c, c.ShouldBindJSON(&id), "Не удалось связать данные с полем id") {
		return
	}
	log.Printf("Удаление записи с ID: %d\n", id.ID)

	if DB := database.GetDB(c); DB != nil {
		obj := DB.Delete(models.Subs{}, id.ID)
		if !utils.CheckError(c, obj.Error, "Не удалось найти запись с таким id") {
			return
		}
	}

	c.JSON(200, gin.H{"Status": "Deleted"})
}