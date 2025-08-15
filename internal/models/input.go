package models

import (
	"github.com/andro-kes/SubAggr/internal/utils"
	"github.com/google/uuid"
)

type Input struct {
	ServiceName string `json:"service_name"`
	Price int `json:"price"`
	UserId uuid.UUID `json:"user_id"`
	StartDate string `json:"start_date"` // MM-YYYY
	EndDate string `json:"end_date"` // MM-YYYY
}

// Создает новую подписку из предоставленных данных
func (input *Input) NewSub() (Subs, error) {
	start, end, err := utils.ParseTime(input.StartDate, input.EndDate)
	if err != nil {
		return Subs{}, err
	}

	return Subs{
		ServiceName: input.ServiceName,
		Price: input.Price,
		UserId: input.UserId,
		StartDate: *start,
		EndDate: end,
	}, nil
}

func (input *Input) IsValid() bool {
	if input.EndDate == "" {
		return utils.IsValidDate(input.StartDate) && utils.IsValidUUID(input.UserId)
	}
	return utils.IsValidDate(input.StartDate) && utils.IsValidUUID(input.UserId) && utils.IsValidDate(input.EndDate)
}