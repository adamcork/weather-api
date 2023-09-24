package main

import (
	handler "github.com/adamcork/weather-api/pkg/api-handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	setup(router)

	port := "localhost:8080"
	router.Run(port)
}

func setup(router *gin.Engine) {
	h := handler.NewAPIHandler()
	router.GET("/weather", h.Weather)
}
