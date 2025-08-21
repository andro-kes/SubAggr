package handlers

import (
	"log/slog"
	"strconv"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/gin-gonic/gin"
)

// ListNotes godoc
// @Summary Получить список подписок
// @Description Возвращает список подписок с возможностью фильтрации
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "Фильтр по ID пользователя"
// @Param service_name query string false "Фильтр по названию сервиса"
// @Param limit query int false "Количество записей (макс 200)"
// @Param offset query int false "Смещение для пагинации"
// @Success 200 {array} models.Subs "Список подписок"
// @Router /SUBS [get]
// @Failure 500 {object} map[string]string "Ошибка сервера"
func ListNotes(c *gin.Context) {
	var notes []models.Subs
	query := database.GetDB(c).Model(&models.Subs{})

	user_id := c.Query("user_id")
	service := c.Query("service_name")
	if user_id != "" {
		query = query.Where("user_id = ?", user_id)
	}
	if service != "" {
		query = query.Where("service_name = ?", service)
	}

	// Pagination with defaults and caps
	const defaultLimit = 50
	const maxLimit = 200
	limit := defaultLimit
	offset := 0
	if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 {
		if v > maxLimit {
			v = maxLimit
		}
		limit = v
	}
	if v, err := strconv.Atoi(c.Query("offset")); err == nil && v >= 0 {
		offset = v
	}

	// Select only required columns for list view
	query = query.Select("id, service_name, price, user_id, start_date, end_date")

	if err := query.Limit(limit).Offset(offset).Find(&notes).Error; err != nil {
		slog.Error("Ошибка выборки списка", slog.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, notes)
}