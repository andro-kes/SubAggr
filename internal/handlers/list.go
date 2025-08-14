package handlers

import (
	"log"

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
// @Success 200 {array} models.Subs "Список подписок"
// @Router /SUBS [get]
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

	query.Find(&notes)

	if len(notes) == 0 {
		log.Println("Ни одной подписки не было найдено")
		c.JSON(200, notes)
		return
	}
	c.JSON(200, notes)
}