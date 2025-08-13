package handlers

import (
	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/utils"
	"github.com/gin-gonic/gin"
)

func CreateNote(c *gin.Context) {
	var sub models.Subs
	if !utils.CheckError(c, c.ShouldBindJSON(&sub), "Невалидные данные для создания записи") {
		return
	}

	if DB := database.GetDB(c); DB != nil {
		obj := DB.Create(&sub)
		if !utils.CheckError(c, obj.Error, "Не удалось создать новую запись") {
			return
		}
	}
	c.JSON(200, gin.H{"Status": "Created"})
}