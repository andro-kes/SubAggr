package app

import "github.com/gin-gonic/gin"

func Run() {
	router := gin.Default()

	router.Run()
}