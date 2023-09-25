package weatherprovider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

func (o *OpenWeatherProvider) GetWarmestDay(lat, long float64) (WeatherResponse, error) {
	values := url.Values{}
	values.Add("lon", fmt.Sprintf("%v", long))
	values.Add("lat", fmt.Sprintf("%v", lat))
	values.Add("units", "metric")
	values.Add("appid", o.apiKey)
	query := values.Encode()

	url := fmt.Sprintf("%s/data/2.5/forecast?%s", o.baseURL, query)
	res, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, errors.New("unable to complete Get request")
	}
	defer res.Body.Close()
	out, err := io.ReadAll(res.Body)
	if err != nil {
		return WeatherResponse{}, errors.New("unable to read response data")
	}

	dataResp := OpenWeatherDataResponse{}
	err = json.Unmarshal(out, &dataResp)
	if err != nil {
		return WeatherResponse{}, err
	}

	if len(dataResp.List) == 0 {
		return WeatherResponse{}, errors.New("no data returned from weather provider")
	}

	return parseWarmestDay(dataResp), nil
}

func parseWarmestDay(r OpenWeatherDataResponse) WeatherResponse {
	resp := WeatherResponse{
		Scale:       "Celcius",
		Temperature: -200.0,
	}

	h := 0

	for _, l := range r.List {
		if float32(l.Main.Temp) > resp.Temperature {
			resp.Temperature = float32(l.Main.Temp)
			resp.WarmestDay = l.DtTxt
			h = l.Main.Humidity
		} else if float32(l.Main.Temp) == resp.Temperature && l.Main.Humidity < h {
			resp.Temperature = float32(l.Main.Temp)
			resp.WarmestDay = l.DtTxt
			h = l.Main.Humidity
		}
	}

	return resp
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
		log.Printf("Error making get request for geo check: %v\n", err)
		return false, errors.New("unable to complete Get request")
	}
	defer res.Body.Close()
	out, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body for geo check.")
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

	log.Printf("%+v\n", geoResp)

	return geoResp[0].Country == "GB", nil
}
