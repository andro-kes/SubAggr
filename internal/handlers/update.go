package handlers

import (
	"log/slog"
	"strconv"
    "time"

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

    if !utils.Ok(updates.IsValid(), "Невалидные данные обновления") {
		c.JSON(400, gin.H{"error": "invalid date"})
		return
	}

	data := make(map[string]interface{})
    // Цена: допускаем 0, различаем отсутствие поля и явное значение
    if updates.Price != nil {
        if *updates.Price < 0 {
            c.JSON(400, gin.H{"error": "price must be >= 0"})
            return
        }
        data["Price"] = *updates.Price
	}
    // Даты: парсим только присланные поля. start_date очищать нельзя.
    var startParsed *time.Time
    var endParsed *time.Time
    if updates.StartDate != nil {
        if *updates.StartDate == "" {
            c.JSON(400, gin.H{"error": "start_date cannot be empty"})
            return
        }
        s, _, err := utils.ParseTime(*updates.StartDate, "")
        if !utils.CheckError(c, err, "Не удалось распарсить start_date") {
            return
        }
        startParsed = s
        data["StartDate"] = *s
	}
    if updates.EndDate != nil {
        if *updates.EndDate == "" {
            // Явный сброс конца подписки (бессрочная)
            data["EndDate"] = nil
        } else {
            _, e, err := utils.ParseTime("01-2000", *updates.EndDate) // фиктивный старт, нам нужен только конец
            if !utils.CheckError(c, err, "Не удалось распарсить end_date") {
                return
            }
            endParsed = e
            data["EndDate"] = endParsed
        }
	}

    // Проверка порядка дат, если обе заданы в запросе (или если одна задана, другая уже есть в БД)
    if DB := database.GetDB(c); DB != nil {
        // Если одна из дат отсутствует в запросе, достанем текущие значения для проверки порядка
        if startParsed == nil || (endParsed == nil && updates.EndDate != nil && *updates.EndDate != "") {
            var current models.Subs
            if err := DB.Select("start_date", "end_date").Where("id = ?", c.Param("id")).First(&current).Error; err == nil {
                if startParsed == nil {
                    startParsed = &current.StartDate
                }
                if endParsed == nil {
                    endParsed = current.EndDate
                }
            }
        }
        if startParsed != nil && endParsed != nil {
            if startParsed.After(*endParsed) {
                c.JSON(400, gin.H{"error": "start_date must be <= end_date"})
                return
            }
        }
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

    c.JSON(200, gin.H{"status": "updated"})
}