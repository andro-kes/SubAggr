package handlers

import (
	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateNote godoc
// @Summary Создать новую подписку
// @Description Добавляет новую запись о подписке пользователя
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param input body models.Input true "Данные подписки"
// @Success 201 {object} models.Subs "Успешно созданная подписка"
// @Failure 400 {object} map[string]string "Невалидные входные данные"
// @Failure 409 {object} map[string]string "Подписка уже существует"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /SUBS [post]
func CreateNote(c *gin.Context) {
	var input models.Input
	if !utils.CheckError(c, c.ShouldBindJSON(&input), "Невалидные данные для создания записи") {
		return
	}

	if !utils.Ok(input.IsValid(), "Невалидная дата") {
		c.JSON(400, gin.H{"error": "invalid date"})
		return
	}

	sub, err := input.NewSub()
	if !utils.CheckError(c, err, "Не удалось создать подписку") {
		return
	}

	if DB := database.GetDB(c); DB != nil {
		if !utils.Ok(isUnique(DB, sub), "Такая подписка уже есть") {
			c.JSON(409, gin.H{"error": "unique constraint violation"})
			return
		}

		obj := DB.Create(&sub)
		if !utils.CheckError(c, obj.Error, "Не удалось создать новую запись") {
			return
		}
	}
	c.JSON(201, sub)
}

func isUnique(db *gorm.DB, sub models.Subs) bool {
	var existingSub models.Subs
	db.Where("user_id = ? AND service_name = ?", sub.UserId, sub.ServiceName).First(&existingSub)
	return existingSub.ID == 0
}