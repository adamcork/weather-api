package weatherprovider

type WeatherResponse struct {
	WarmestDay  string
	Temperature float32
	Scale       string
}

type OpenWeatherGeoResponse []struct {
	Name       string            `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Country    string            `json:"country"`
	State      string            `json:"state"`
}
