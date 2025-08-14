package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Aggr interface {
	IsValidUUID() bool
	IsValidDate() bool
}

type Subs struct {
	gorm.Model
	ServiceName string `json:"service_name"`
	Price int `json:"price"`
	UserId uuid.UUID `json:"user_id"`
	StartDate string `json:"start_date"` // MM-YYYY
	EndDate string `json:"end_date"` // MM-YYYY
}

func (sub *Subs) IsValidUUID() bool {
    _, err := uuid.Parse(sub.UserId.String())
    return err == nil
}

func (sub *Subs) IsValidDate() bool {
    _, errStart := time.Parse("01-2006", sub.StartDate)
	if sub.EndDate == "" {
		return errStart == nil
	}
	_, errEnd := time.Parse("01-2006", sub.EndDate)
    return errStart == nil && errEnd == nil
}

type Filters struct {
	ServiceName string `json:"service_name"`
	UserId uuid.UUID `json:"user_id"`
	StartDate string `json:"start_date"` // MM-YYYY
	EndDate string `json:"end_date"` // MM-YYYY
}

func (sub *Filters) IsValidUUID() bool {
    _, err := uuid.Parse(sub.UserId.String())
    return err == nil
}

func (sub *Filters) IsValidDate() bool {
    _, errStart := time.Parse("01-2006", sub.StartDate)
	if sub.EndDate == "" {
		return errStart == nil
	}
	_, errEnd := time.Parse("01-2006", sub.EndDate)
    return errStart == nil && errEnd == nil
}

func IsValid(obj Aggr) bool {
	return obj.IsValidDate() && obj.IsValidUUID()
}