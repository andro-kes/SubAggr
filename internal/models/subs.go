package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subs struct {
	gorm.Model
	ServiceName string `json:"service_name" gorm:"index;uniqueIndex:uniq_user_service"`
	Price int `json:"price"`
	UserId uuid.UUID `json:"user_id" gorm:"index;uniqueIndex:uniq_user_service"`
	StartDate time.Time `json:"start_date" gorm:"index"` // MM-YYYY
	EndDate *time.Time `json:"end_date" gorm:"index"` // MM-YYYY
}