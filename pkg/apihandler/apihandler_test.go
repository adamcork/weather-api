package apihandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
		resp: "Monday",
		err:  nil,
	}

	router := setupRouter(m)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weather?long=123.456&lat=234.567", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := `"Monday"`
	assert.JSONEq(t, expected, w.Body.String())

}

type mockWeatherProvider struct {
	resp string
	err  error
}

func (m *mockWeatherProvider) GetWarmestDay(lat, long float64) (string, error) {
	return m.resp, m.err
}
