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
	CheckUKLocation(lat, long float64) (bool, error)
}

type HistoryService interface {
	SaveRequest(lat, long float64)
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

	lat, err := strconv.ParseFloat(lt, 32)
	if err != nil {
		errs = append(errs, "lat parameter could not be parsed as a float.")
	}

	long, err := strconv.ParseFloat(lg, 32)
	if err != nil {
		errs = append(errs, "long parameter could not be parsed as a float.")
	}

	if len(errs) > 0 {
		resp := WeatherResponse{
			Errors: errs,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	pattern := `^[-+]?[0-9]*\.?[0-9]{0,6}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(lt) {
		errs = append(errs, "lat parameter has too many decimal places.")
	}

	if !regex.MatchString(lg) {
		errs = append(errs, "long parameter has too many decimal places.")
	}

	if len(errs) > 0 {
		resp := WeatherResponse{
			Errors: errs,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	isUK, _ := h.provider.CheckUKLocation(lat, long)
	if !isUK {
		errs = append(errs, "Only UK locations are permitted.")
		resp := WeatherResponse{
			Errors: errs,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	providerResp, err := h.provider.GetWarmestDay(lat, long)
	if err != nil {
		errs = append(errs, err.Error())
		resp := WeatherResponse{
			Errors: errs,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := WeatherResponse{
		WarmestDay: providerResp.WarmestDay,
		Temperature: Temperature{
			Value: providerResp.Temperature,
			Scale: providerResp.Scale,
		},
	}

	c.JSON(http.StatusOK, resp)
}
