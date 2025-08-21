package utils

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MustNotError(err error, msg string) {
	if err != nil {
		slog.Error("fatal error", slog.String("msg", msg), slog.String("error", err.Error()))
		panic(err)
	}
}

func Ok(ok bool, msg string) bool {
	if !ok {
		slog.Warn(msg)
		return false
	}
	return true
}

/*
Проверяет, есть ли ошибка при работе хэндлера.
true - ошибки нет,
false - есть ошибка
*/
func CheckError(c *gin.Context, err error, msg string) bool {
	if err != nil {
		slog.Warn(msg, slog.String("error", err.Error()))
		status := 400
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = 404
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return false
	}
	return true
}