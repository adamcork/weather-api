package apihandler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/adamcork/weather-api/pkg/weatherprovider"
	"github.com/gin-gonic/gin"
)

type WeatherProvider interface {
	GetWarmestDay(lat, long float64) (weatherprovider.WeatherResponse, error)
}

type APIHandler struct {
	provider WeatherProvider
}

type WeatherResponse struct {
	WarmestDay  string      `json:"warmest-day"`
	Temperature Temperature `json:"temperature"`
	Errors      []string    `json:"errors,omitempty"`
}

type Temperature struct {
	Value float32 `json:"value"`
	Scale string  `json:"scale"` // Celcius, Farenheit (Kelvin???)
}

func NewAPIHandler(provider WeatherProvider) *APIHandler {
	return &APIHandler{
		provider: provider,
	}
}

func (h *APIHandler) Weather(c *gin.Context) {
	errs := []string{}
	lt := c.Query("lat")
	lg := c.Query("long")

	pattern := `^[-+]?[0-9]*\.?[0-9]{0,6}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(lt) {
		errs = append(errs, "lat parameter has too many decimal places.")
		resp := WeatherResponse{
			Errors: errs,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	lat, _ := strconv.ParseFloat(lt, 32)
	long, _ := strconv.ParseFloat(lg, 32)

	providerResp, _ := h.provider.GetWarmestDay(lat, long)

	resp := WeatherResponse{
		WarmestDay: providerResp.WarmestDay,
		Temperature: Temperature{
			Value: providerResp.Temperature,
			Scale: providerResp.Scale,
		},
	}

	c.JSON(http.StatusOK, resp)
}
