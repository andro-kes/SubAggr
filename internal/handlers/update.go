package handlers

import (
	"log"
	"strconv"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/internal/utils"
	"github.com/gin-gonic/gin"
)

// UpdateNote godoc
// @Summary Обновить подписку
// @Description Обновляет существующую подписку по ID
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param input body models.Subs true "Обновленные данные подписки"
// @Success 200 {object} models.Subs "Обновленная подписка"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /SUBS/{id} [put]
func UpdateNote(c *gin.Context) {
	var updatedSub models.Subs
	if !utils.CheckError(c, c.ShouldBindJSON(&updatedSub), "Невалидные данные для связывания записи") {
		return
	}
	log.Printf("Обновление записи с ID: %d\n", updatedSub.ID)

	if !utils.IsValid(models.IsValid(&updatedSub), "Невалидная дата") {
		c.JSON(400, gin.H{"Error": "Невалидная дата"})
		return
	}

	ID := c.Param("id")
	log.Printf("Чтение записи с ID: %s\n", ID)

	id, err := strconv.ParseUint(ID, 10, 64)
	if !utils.CheckError(c, err, "Невалидный id") {
		return
	}

	if DB := database.GetDB(c); DB != nil {
		obj := DB.Where("id = ?", id).Updates(gin.H{
			"price": updatedSub.Price,
			"start_date": updatedSub.StartDate,
			"end_date": updatedSub.EndDate,
		})
		if !utils.CheckError(c, obj.Error, "Не удалось обновить данные записи") {
			return
		}
	}

	c.JSON(200, updatedSub)
}