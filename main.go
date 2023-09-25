package main

import (
	"github.com/adamcork/weather-api/pkg/apihandler"
	"github.com/adamcork/weather-api/pkg/history"
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
	p := weatherprovider.NewOpenWeatherProvider("http://api.openweathermap.org", "apiKey")
	hs := history.NewFSHistoryService()
	h := apihandler.NewAPIHandler(p, hs)
	router.GET("/weather", h.Weather)
	// TODO: Add endpoint to map history handler
}
