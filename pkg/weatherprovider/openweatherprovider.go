package weatherprovider

type OpenWeatherProvider struct {
}

func NewOpenWeatherProvider() *OpenWeatherProvider {
	return &OpenWeatherProvider{}
}

func (o *OpenWeatherProvider) GetWarmestDay(lat, long float64) (string, error) {
	// TODO: Implement
	return "not implemented", nil
}
