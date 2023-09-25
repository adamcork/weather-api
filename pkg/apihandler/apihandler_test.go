package apihandler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamcork/weather-api/pkg/history"
	"github.com/adamcork/weather-api/pkg/weatherprovider"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(m *mockWeatherProvider, hs *mockHistoryService) *gin.Engine {
	r := gin.Default()
	sut := NewAPIHandler(m, hs)
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
		wamestErr: nil,
		isUK:      true,
	}

	router := setupRouter(m, &mockHistoryService{})

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
		resp:      weatherprovider.WeatherResponse{},
		wamestErr: nil,
	}

	router := setupRouter(m, &mockHistoryService{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather?long=123.456&lat=234.5672154", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"warmest-day": "", "temperature": {"value": 0, "scale": ""}, "errors": ["lat parameter has too many decimal places."]}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLongTooManyDP(t *testing.T) {
	// Not expecting call to weather provider - check decimal places in handler.
	// If implemented in weather provider, would have to be re-implemented in all providers.
	m := &mockWeatherProvider{
		resp:      weatherprovider.WeatherResponse{},
		wamestErr: nil,
	}

	router := setupRouter(m, &mockHistoryService{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather?long=123.4567895&lat=234.567215", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"warmest-day": "", "temperature": {"value": 0, "scale": ""}, "errors": ["long parameter has too many decimal places."]}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestBothTooManyDP(t *testing.T) {
	// Not expecting call to weather provider - check decimal places in handler.
	// If implemented in weather provider, would have to be re-implemented in all providers.
	m := &mockWeatherProvider{
		resp:      weatherprovider.WeatherResponse{},
		wamestErr: nil,
		isUK:      true,
	}

	router := setupRouter(m, &mockHistoryService{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather?long=123.4567895&lat=234.56721564", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"warmest-day": "", "temperature": {"value": 0, "scale": ""}, "errors": ["lat parameter has too many decimal places.","long parameter has too many decimal places."]}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLongNonNumeric(t *testing.T) {
	// Not expecting call to weather provider - check decimal places in handler.
	// If implemented in weather provider, would have to be re-implemented in all providers.
	m := &mockWeatherProvider{
		resp:      weatherprovider.WeatherResponse{},
		wamestErr: nil,
	}

	router := setupRouter(m, &mockHistoryService{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather?long=fail&lat=234.567215", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"warmest-day": "", "temperature": {"value": 0, "scale": ""}, "errors": ["long parameter could not be parsed as a float."]}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestOutsideUK(t *testing.T) {
	m := &mockWeatherProvider{
		isUK: false,
	}

	router := setupRouter(m, &mockHistoryService{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather?long=123.456&lat=234.567", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"warmest-day": "", "temperature": {"value": 0, "scale": ""}, "errors": ["Only UK locations are permitted."]}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestWeatherResponseError(t *testing.T) {
	e := "Error from weather provider."
	m := &mockWeatherProvider{
		resp:      weatherprovider.WeatherResponse{},
		wamestErr: errors.New(e),
		isUK:      true,
	}

	router := setupRouter(m, &mockHistoryService{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather?long=123.456&lat=234.567", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expected := fmt.Sprintf(`{"warmest-day": "", "temperature": {"value": 0, "scale": ""}, "errors": ["%s"]}`, e)
	assert.JSONEq(t, expected, w.Body.String())
}

type mockWeatherProvider struct {
	resp         weatherprovider.WeatherResponse
	wamestErr    error
	isUK         bool
	uKCheckError error
}

func (m *mockWeatherProvider) GetWarmestDay(lat, long float64) (weatherprovider.WeatherResponse, error) {
	return m.resp, m.wamestErr
}

func (m *mockWeatherProvider) CheckUKLocation(lat, long float64) (bool, error) {
	return m.isUK, m.uKCheckError
}

type mockHistoryService struct {
}

func (h *mockHistoryService) SaveRequest(history.WeatherRequest) error {
	// TODO: Implement save weather request tests
	return nil
}

func (h *mockHistoryService) GetHistory(orderBy string, limit int) ([]history.WeatherRequest, error) {
	// TODO: Implement gethistory handler tests
	return nil, nil
}
