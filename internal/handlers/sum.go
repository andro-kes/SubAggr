package handlers

import (
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
// @Router /SUBS/SUMMARY [post]
func SumPriceSubs(c *gin.Context) {
	var filters models.Filters
	if !utils.CheckError(c, c.ShouldBindJSON(&filters), "Не удалось связать фильтры") {
		return
	}

	if !utils.Ok(filters.IsValid(), "Невалидная дата") {
		c.JSON(400, gin.H{"error": "invalid date"})
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

	total, err := aggregateTotal(DB, sub)
	if !utils.CheckError(c, err, "Ошибка агрегации суммы") {
		return
	}

	if total == 0 {
		c.JSON(200, gin.H{"message": "no active subscriptions"})
		return
	}
	c.JSON(200, gin.H{"total": total})
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

// Выполняет агрегацию суммы в БД
func aggregateTotal(db *gorm.DB, filters models.Subs) (int64, error) {
	var total int64
	query := `
WITH params AS (
    SELECT date_trunc('month', $1::date) AS start_month,
           date_trunc('month', $2::date) AS end_month
), eligible AS (
    SELECT s.price,
           GREATEST(date_trunc('month', s.start_date), p.start_month) AS from_month,
           LEAST(date_trunc('month', COALESCE(s.end_date, p.end_month)), p.end_month) AS to_month
    FROM subs s
    CROSS JOIN params p
    WHERE s.start_date <= p.end_month
      AND (s.end_date IS NULL OR s.end_date >= p.start_month)
      AND ($3::uuid IS NULL OR s.user_id = $3::uuid)
      AND ($4::text IS NULL OR $4 = '' OR s.service_name = $4)
), months AS (
    SELECT price,
           COUNT(*)::int AS m
    FROM eligible e
    JOIN generate_series(
        e.from_month,
        e.to_month - interval '1 month',
        interval '1 month'
    ) gs ON TRUE
    WHERE e.to_month >= e.from_month
    GROUP BY price
)
SELECT COALESCE(SUM(price * m), 0) AS total
FROM months;`

	var userParam interface{}
	if filters.UserId != uuid.Nil {
		userParam = filters.UserId
	} else {
		userParam = nil
	}
	var serviceParam interface{}
	if filters.ServiceName != "" {
		serviceParam = filters.ServiceName
	} else {
		serviceParam = nil
	}

	err := db.Raw(query, filters.StartDate, *filters.EndDate, userParam, serviceParam).Scan(&total).Error
	return total, err
}
