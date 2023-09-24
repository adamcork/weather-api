package apihandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WeatherProvider interface {
	GetWarmestDay(lat, long float64) (string, error)
}

type APIHandler struct {
	provider WeatherProvider
}

func NewAPIHandler(provider WeatherProvider) *APIHandler {
	return &APIHandler{
		provider: provider,
	}
}

func (h *APIHandler) Weather(c *gin.Context) {
	lt := c.Query("lat")
	lg := c.Query("long")
	lat, _ := strconv.ParseFloat(lt, 32)
	long, _ := strconv.ParseFloat(lg, 32)

	day, _ := h.provider.GetWarmestDay(lat, long)
	c.JSON(http.StatusOK, day)
}
