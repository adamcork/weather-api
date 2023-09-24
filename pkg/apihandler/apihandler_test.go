package apihandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamcork/weather-api/pkg/weatherprovider"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(m *mockWeatherProvider) *gin.Engine {
	r := gin.Default()
	sut := NewAPIHandler(m)
	r.GET("/weather", sut.Weather)
	return r
}

func TestGetWeather(t *testing.T) {
	m := &mockWeatherProvider{
		resp: weatherprovider.WeatherResponse{
			WarmestDay:  "Monday",
			Temperature: 15.3,
			Scale:       "Celcius",
		},
		err: nil,
	}

	router := setupRouter(m)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather?long=123.456&lat=234.567", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := `{"warmest-day": "Monday", "temperature": {"value": 15.3, "scale": "Celcius"}}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLatTooManyDP(t *testing.T) {
	// Not expecting call to weather provider - check decimal places in handler.
	// If implemented in weather provider, would have to be re-implemented in all providers.
	m := &mockWeatherProvider{
		resp: weatherprovider.WeatherResponse{},
		err:  nil,
	}

	router := setupRouter(m)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather?long=123.456&lat=234.5672154", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"warmest-day": "", "temperature": {"value": 0, "scale": ""}, "errors": ["lat parameter has too many decimal places."]}`
	assert.JSONEq(t, expected, w.Body.String())
}

type mockWeatherProvider struct {
	resp weatherprovider.WeatherResponse
	err  error
}

func (m *mockWeatherProvider) GetWarmestDay(lat, long float64) (weatherprovider.WeatherResponse, error) {
	return m.resp, m.err
}
