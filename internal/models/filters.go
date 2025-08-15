package models

import (
	"github.com/andro-kes/SubAggr/internal/utils"
	"github.com/google/uuid"
)

type Filters struct {
	ServiceName string `json:"service_name"`
	UserId uuid.UUID `json:"user_id"`
	StartDate string `json:"start_date"` // MM-YYYY
	EndDate string `json:"end_date"` // MM-YYYY
}

func (filters *Filters) NewSub() (Subs, error) {
	start, end, err := utils.ParseTime(filters.StartDate, filters.EndDate)
	if err != nil {
		return Subs{}, err
	}

	return Subs{
		ServiceName: filters.ServiceName,
		UserId: filters.UserId,
		StartDate: *start,
		EndDate: end,
	}, nil
}

func (filters *Filters) IsValid() bool {
	return utils.IsValidDate(filters.StartDate) && utils.IsValidUUID(filters.UserId) && utils.IsValidDate(filters.EndDate)
}