package history

import "time"

type WeatherRequest struct {
	Timestamp   time.Time
	Long        float64
	Lat         float64
	Success     bool
	WarmestDay  string
	Temperature float64
	Scale       string
}
