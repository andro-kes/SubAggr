package handlers

import (
	"log/slog"
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
// @Param input body models.Updates true "Обновленные данные подписки"
// @Success 200 {object} models.Subs "Обновленная подписка"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /SUBS/{id} [put]
func UpdateNote(c *gin.Context) {
	var updates models.Updates

	if !utils.CheckError(c, c.ShouldBindJSON(&updates), "Невалидные данные для связывания записи") {
		return
	}

	if !utils.Ok(updates.IsValid(), "Невалидная дата") {
		c.JSON(400, gin.H{"error": "invalid date"})
		return
	}

	sub, err := updates.NewSub()
	if !utils.CheckError(c, err, "Не удалось организовать данные для обновления подписки") {
		return
	}

	data := make(map[string]interface{})
	if sub.Price != 0 {
		data["Price"] = sub.Price
	}
	if !sub.StartDate.IsZero() {
		data["StartDate"] = sub.StartDate
	}
	if sub.EndDate != nil {
		data["EndDate"] = sub.EndDate
	}

	if len(data) == 0 {
		c.JSON(200, gin.H{"message": "no changes"})
		return
	}

	ID := c.Param("id")
	slog.Debug("Обновление записи", slog.String("id", ID))

	id, err := strconv.ParseUint(ID, 10, 64)
	if !utils.CheckError(c, err, "Невалидный id") {
		return
	}

	if DB := database.GetDB(c); DB != nil {
		obj := DB.Model(models.Subs{}).Where("id = ?", id).Updates(data)
		if !utils.CheckError(c, obj.Error, "Не удалось обновить данные записи") {
			return
		}
		if obj.RowsAffected == 0 {
			c.JSON(404, gin.H{"error": "subscription not found"})
			return
		}
	}

	c.JSON(200, updates)
}