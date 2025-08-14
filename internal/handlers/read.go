package handlers

import (
	"log"
	"strconv"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/internal/utils"
	"github.com/gin-gonic/gin"
)

// ReadNote godoc
// @Summary Получить подписку
// @Description Возвращает подписку по указанному ID
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} models.Subs "Данные подписки"
// @Failure 400 {object} map[string]string "Невалидный ID"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /SUBS/{id} [get]
func ReadNote(c *gin.Context) {
	ID := c.Param("id")
	log.Printf("Чтение записи с ID: %s\n", ID)

	id, err := strconv.ParseUint(ID, 10, 64)
	if !utils.CheckError(c, err, "Невалидный id") {
		return
	}

	var sub models.Subs
	if DB := database.GetDB(c); DB != nil {
		obj := DB.Where("id = ?", uint(id)).First(&sub)
		if !utils.CheckError(c, obj.Error, "Не удалось найти запись с таким ID") {
			return
		}
	}

	c.JSON(200, sub)
}