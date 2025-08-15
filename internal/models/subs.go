package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subs struct {
	gorm.Model
	ServiceName string `json:"service_name"`
	Price int `json:"price"`
	UserId uuid.UUID `json:"user_id"`
	StartDate time.Time `json:"start_date"` // MM-YYYY
	EndDate *time.Time `json:"end_date"` // MM-YYYY
}