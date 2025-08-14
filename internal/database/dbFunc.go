package database

import (
	"github.com/andro-kes/SubAggr/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDB(c *gin.Context) *gorm.DB {
	db, ok := c.Get("DB")
	if !utils.IsValid(ok, "Контекст не содержит базу данных") {
		return nil
	}
	DB, ok := db.(*gorm.DB)
	if !utils.IsValid(ok, "Невалидный тип данных для БД") {
		return nil
	}
	return DB
}