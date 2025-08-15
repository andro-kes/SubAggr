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
// @Param input body models.Updates true "Обновленные данные подписки"
// @Param input body models.Subs true "Обновленные данные подписки"
// @Success 200 {object} models.Subs "Обновленная подписка"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /SUBS/{id} [put]
func UpdateNote(c *gin.Context) {
	var updates models.Updates

	if !utils.CheckError(c, c.ShouldBindJSON(&updates), "Невалидные данные для связывания записи") {
		return
	}

	if !utils.Ok(updates.IsValid(), "Невалидная дата") {
		c.JSON(400, gin.H{"Error": "Невалидная дата"})
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
		c.JSON(200, gin.H{"Message": "No changes"})
	}

	ID := c.Param("id")
	log.Printf("Обновление записи с ID: %s\n", ID)

	id, err := strconv.ParseUint(ID, 10, 64)
	if !utils.CheckError(c, err, "Невалидный id") {
		return
	}

	if DB := database.GetDB(c); DB != nil {
		obj := DB.Model(models.Subs{}).Where("id = ?", id).Updates(data)
		if !utils.CheckError(c, obj.Error, "Не удалось обновить данные записи") {
			return
		}
	}

	c.JSON(200, updates)
}