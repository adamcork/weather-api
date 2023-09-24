package main

import (
	"github.com/adamcork/weather-api/pkg/apihandler"
	"github.com/adamcork/weather-api/pkg/weatherprovider"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	setup(router)

	port := "localhost:8080"
	router.Run(port)
}

func setup(router *gin.Engine) {
	// TODO: add config
	p := weatherprovider.NewOpenWeatherProvider("", "")
	h := apihandler.NewAPIHandler(p)
	router.GET("/weather", h.Weather)
}
