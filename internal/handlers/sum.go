package handlers

import (
	"log"
	"time"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SumPriceSubs godoc
// @Summary Сумма подписок
// @Description Возвращает сумму цен подписок за период с фильтрацией
// @Tags Aggregations
// @Accept json
// @Produce json
// @Param input body models.Filters true "Параметры фильтрации" example:{"service_name":"Yandex Plus","user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025","end_date":"12-2025"}
// @Success 200 {object} map[string]interface{} "Сумма подписок или сообщение"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /SUMSUBS [post]
func SumPriceSubs(c *gin.Context) {
	var filters models.Filters
	if !utils.CheckError(c, c.ShouldBindJSON(&filters), "Не удалось связать фильтры") {
		return
	}

	if !utils.Ok(filters.IsValid(), "Невалидная дата") {
		c.JSON(400, gin.H{"Error": "Невалидная дата"})
		return
	}

	DB := database.GetDB(c)
	if DB == nil {
		return
	}

	sub, err := filters.NewSub()
	if !utils.CheckError(c, err, "Не удалось организовать данные для фильтрации") {
		return
	}

	filterSubs := filter(DB, sub)
	log.Printf("Найдено активных подписок: %d\n", len(filterSubs))
	sum := summary(filterSubs, sub)

	if sum == 0{
		c.JSON(200, gin.H{"Message": "Нет активных подписок"})
		return
	}
	c.JSON(200, gin.H{"Сумма подписок": sum})
}

// Находит активные подписки, учитывая фильтры
func filter(db *gorm.DB, filters models.Subs) []models.Subs {
	var sumSubs []models.Subs
	query := db

	if filters.UserId != uuid.Nil {
		query = query.Where("user_id = ?", filters.UserId)
	}
	if filters.ServiceName != "" {
		query = query.Where("service_name = ?", filters.ServiceName)
	}
	query = query.Where("start_date <= ?", filters.EndDate)
	query = query.Where("(end_date >= ? OR end_date IS NULL)", filters.StartDate)

	// Select only what we need for pricing
	query = query.Select("price, start_date, end_date")

	query.Find(&sumSubs)

	return sumSubs
}

// Возвращает стоимость подписки за выбранный период
func summary(subs []models.Subs, filters models.Subs) int {
	sum := 0
	var start, end time.Time
	for _, sub := range subs {
		if sub.StartDate.After(filters.StartDate) {
			start = sub.StartDate
		} else {
			start = filters.StartDate
		}

		if sub.EndDate == nil {
			end = *filters.EndDate
		} else if sub.EndDate.Before(*filters.EndDate) {
			end = *sub.EndDate
		} else {
			end = *filters.EndDate
		}

		sum += calculatePrice(sub.Price, start, end)
	}
	return sum
}

func calculatePrice(price int, start, end time.Time) int {
	months := monthsBetween(start, end)
	return months * price
}

// monthsBetween returns the number of whole months between start (inclusive) and end (exclusive),
// matching previous loop semantics that incremented one month at a time until start == end.
func monthsBetween(start, end time.Time) int {
	if !start.Before(end) {
		return 0
	}
	y1, m1, _ := start.Date()
	y2, m2, _ := end.Date()
	return int((y2-y1)*12) + int(m2-m1)
}