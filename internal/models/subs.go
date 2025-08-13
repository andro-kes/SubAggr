package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subs struct {
	gorm.Model
	ServiceName string `json:"service_name"`
	Price int `json:"price"`
	UserId uuid.UUID `json:"user_id"`
	StartDate string `json:"start_date"`
}