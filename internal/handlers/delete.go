package handlers

import (
	"log/slog"
	"strconv"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/internal/utils"
	"github.com/gin-gonic/gin"
)

// DeleteNote godoc
// @Summary Удалить подписку
// @Description Удаляет существующую подписку по ID
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} map[string]string "Успешное удаление"
// @Failure 400 {object} map[string]string "Невалидный ID"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /SUBS/{id} [delete]
func DeleteNote(c *gin.Context) {
	ID := c.Param("id")
	slog.Debug("Удаление записи", slog.String("id", ID))

	id, err := strconv.ParseUint(ID, 10, 64)
	if !utils.CheckError(c, err, "Невалидный id") {
		return
	}

	if DB := database.GetDB(c); DB != nil {
		obj := DB.Delete(&models.Subs{}, id)
		if !utils.CheckError(c, obj.Error, "Не удалось удалить запись") {
			return
		}
		if obj.RowsAffected == 0 {
			c.JSON(404, gin.H{"error": "subscription not found"})
			return
		}
	}

	c.JSON(200, gin.H{"status": "deleted"})
}