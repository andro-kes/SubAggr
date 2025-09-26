package models

import (
	"github.com/andro-kes/SubAggr/internal/utils"
)

type Updates struct {
    // Используем указатели, чтобы различать "поле не прислано" и "значение по умолчанию"
    Price *int `json:"price"`
    StartDate *string `json:"start_date"` // MM-YYYY, nil — не менять, "" — сбросить?
    EndDate *string `json:"end_date"` // MM-YYYY, nil — не менять, "" — сбросить (NULL)
}

// В контексте обновлений парсинг выполняется в хэндлере, где известно какие поля реально присланы
// поэтому здесь вспомогательная функция не нужна.

func (updates *Updates) IsValid() bool {
    // Валидируем только те поля, что реально присланы (не nil)
    if updates.StartDate != nil && *updates.StartDate != "" {
        if !utils.IsValidDate(*updates.StartDate) {
            return false
        }
    }
    if updates.EndDate != nil && *updates.EndDate != "" {
        if !utils.IsValidDate(*updates.EndDate) {
            return false
        }
    }
    // Цена может быть нулевой; если прислана — только проверяем, что не отрицательная
    if updates.Price != nil && *updates.Price < 0 {
        return false
    }
    return true
}