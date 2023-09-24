package weatherprovider

type OpenWeatherProvider struct {
}

func NewOpenWeatherProvider() *OpenWeatherProvider {
	return &OpenWeatherProvider{}
}

func (o *OpenWeatherProvider) GetWarmestDay(lat, long float64) (WeatherResponse, error) {
	// TODO: Implement
	return WeatherResponse{}, nil
}

func (o *OpenWeatherProvider) CheckUKLocation(lat, long float64) (bool, error) {
	// TODO: Implement
	return true, nil
}
