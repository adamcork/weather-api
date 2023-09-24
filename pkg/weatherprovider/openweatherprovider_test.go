package weatherprovider

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type geoTest struct {
	name     string
	location string
	lat      float64
	long     float64
	apiKey   string
	expResp  bool
}

func TestGeoCheck(t *testing.T) {
	tests := []geoTest{
		{
			name:     "paris",
			location: "paris",
			lat:      12.54,
			long:     58.21,
			apiKey:   "apikey123456",
			expResp:  false,
		},
		{
			name:     "belfast",
			location: "belfast",
			lat:      23.98,
			long:     71.65,
			apiKey:   "apikey665544",
			expResp:  true,
		},
		{
			name:     "brighton",
			location: "brighton",
			lat:      66.66,
			long:     1.44,
			apiKey:   "apikey774411",
			expResp:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geoResp := getSample(tt.location)

			testHandler := func(w http.ResponseWriter, r *http.Request) {
				lt := r.URL.Query().Get("lat")
				lg := r.URL.Query().Get("lon")
				limit := r.URL.Query().Get("limit")
				key := r.URL.Query().Get("appid")

				assert.Equal(t, fmt.Sprintf("%v", tt.lat), lt)
				assert.Equal(t, fmt.Sprintf("%v", tt.long), lg)
				assert.Equal(t, "5", limit)
				assert.Equal(t, tt.apiKey, key)
				fmt.Fprint(w, geoResp)
			}

			svr := httptest.NewServer(http.HandlerFunc(testHandler))
			defer svr.Close()

			sut := NewOpenWeatherProvider(svr.URL, tt.apiKey)

			b, err := sut.CheckUKLocation(tt.lat, tt.long)
			assert.Nil(t, err)
			assert.Equal(t, tt.expResp, b)
		})
	}

}

func getSample(s string) string {
	content, err := os.ReadFile(fmt.Sprintf("testdata/openweather_geo_%s.json", s))
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}
