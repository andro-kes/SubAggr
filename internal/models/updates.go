package models

import (
	"github.com/andro-kes/SubAggr/internal/utils"
)

type Updates struct {
	Price int `json:"price"`
	StartDate string `json:"start_date"` // MM-YYYY
	EndDate string `json:"end_date"` // MM-YYYY
}

func (updates *Updates) NewSub() (Subs, error) {
	start, end, err := utils.ParseTime(updates.StartDate, updates.EndDate)
	if err != nil {
		return Subs{}, err
	}

	return Subs{
		Price: updates.Price,
		StartDate: *start,
		EndDate: end,
	}, nil
}

func (updates *Updates) IsValid() bool {
	if updates.EndDate == "" && updates.StartDate == "" {
		return true
	} else if updates.EndDate == "" {
		return utils.IsValidDate(updates.StartDate)
	} else if updates.StartDate == "" {
		return utils.IsValidDate(updates.EndDate)
	}
	return utils.IsValidDate(updates.StartDate) && utils.IsValidDate(updates.EndDate)
}