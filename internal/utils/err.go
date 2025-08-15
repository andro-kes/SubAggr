package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func MustNotError(err error, msg string) {
	if err != nil {
		log.Fatalf("Message: %s\nError: %s\n", msg, err.Error())
	}
}

func Ok(ok bool, msg string) bool {
	if !ok {
		log.Println(msg)
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
		log.Println(msg)
		c.JSON(400, gin.H{"Error": err.Error()})
		return false
	}
	return true
}