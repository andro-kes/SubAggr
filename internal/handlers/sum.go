package handlers

import (
	"time"

	"github.com/andro-kes/SubAggr/internal/database"
	"github.com/andro-kes/SubAggr/internal/models"
	"github.com/andro-kes/SubAggr/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SumPriceSubs godoc
// @Summary Сумма подписок
// @Description Возвращает сумму цен подписок за период с фильтрацией
// @Tags Aggregations
// @Accept json
// @Produce json
// @Param input body models.Filters true "Параметры фильтрации"
// @Success 200 {object} map[string]interface{} "Сумма подписок или сообщение"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /SUMSUBS [post]
func SumPriceSubs(c *gin.Context) {
	var filters models.Filters
	if !utils.CheckError(c, c.ShouldBindJSON(&filters), "Не удалось связать фильтры") {
		return
	}

	if !utils.IsValid(models.IsValid(&filters), "Невалидная дата") {
		c.JSON(400, gin.H{"Error": "Невалидная дата"})
		return
	}

	DB := database.GetDB(c)
	if DB == nil {
		return
	}

	filterSubs := filter(DB, filters)
	sum := summary(filterSubs)

	if sum == 0{
		c.JSON(200, gin.H{"Message": "Нет активных подписок"})
		return
	}
	c.JSON(200, gin.H{"Сумма подписок": sum})
}

func filter(db *gorm.DB, filters models.Filters) []models.Subs {
	var sumSubs []models.Subs

	start, end, err := parseTime(filters.StartDate, filters.EndDate)
	if err != nil {
		return sumSubs
	}

    db.Where(
        "user_id = ? AND service_name = ? AND start_date <= ? AND (end_date >= ? OR end_date IS NULL)",
        filters.UserId,
        filters.ServiceName,
        end,
        start,
    ).Find(&sumSubs)

	return sumSubs
}

func parseTime(startDate, endDate string) (time.Time, time.Time, error) {
	start, err := time.Parse("01-2006", startDate)
	if !utils.IsValid(err == nil, "Не удалось распарсить дату старта при фильтрации") {
		return time.Now(), time.Now(), err
	}
	end, err := time.Parse("01-2006", endDate)
	if !utils.IsValid(err == nil, "Не удалось распарсить дату окончания при фильтрации") {
		return time.Now(), time.Now(), err
	}
	return start, end, nil
}

func summary(subs []models.Subs) int {
	sum := 0
	for _, sub := range subs {
		start, end, err := parseTime(sub.StartDate, sub.EndDate)
		if err != nil {
			return sum
		}
		sum += calculatePrice(sub.Price, start, end)
	}
	return sum
}

const NUMBER_OF_HOURS_PER_MONTH = 720

func calculatePrice(price int, start, end time.Time) int {
	duration := end.Sub(start)
	commonPrice := int(duration.Hours())
	return (commonPrice / NUMBER_OF_HOURS_PER_MONTH) * price
}