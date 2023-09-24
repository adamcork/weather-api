package weatherprovider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type OpenWeatherProvider struct {
	baseURL string
	apiKey  string
}

func NewOpenWeatherProvider(baseUrl, apiKey string) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		baseURL: baseUrl,
		apiKey:  apiKey,
	}
}

// api.openweathermap.org/data/2.5/forecast?lat=50.8226&lon=-0.152778&units=metric&appid={apikey}
// api.openweathermap.org/geo/1.0/reverse?lat=54.5973&lon=-5.9301&limit=5&appid={apikey}

func (o *OpenWeatherProvider) GetWarmestDay(lat, long float64) (WeatherResponse, error) {
	// TODO: Implement
	return WeatherResponse{}, nil
}

func (o *OpenWeatherProvider) CheckUKLocation(lat, long float64) (bool, error) {
	values := url.Values{}
	values.Add("lon", fmt.Sprintf("%v", long))
	values.Add("lat", fmt.Sprintf("%v", lat))
	values.Add("limit", "5")
	values.Add("appid", o.apiKey)
	query := values.Encode()

	url := fmt.Sprintf("%s/geo/1.0/reverse?%s", o.baseURL, query)
	res, err := http.Get(url)
	if err != nil {
		return false, errors.New("unable to complete Get request")
	}
	defer res.Body.Close()
	out, err := io.ReadAll(res.Body)
	if err != nil {
		return false, errors.New("unable to read response data")
	}

	geoResp := OpenWeatherGeoResponse{}
	err = json.Unmarshal(out, &geoResp)
	if err != nil {
		return false, err
	}

	if len(geoResp) == 0 {
		return false, nil
	}

	return geoResp[0].Country == "GB", err
}
