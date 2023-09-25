package history

type FSHistoryService struct {
	requests []WeatherRequest
}

func NewFSHistoryService() *FSHistoryService {
	return &FSHistoryService{}
}

func (h *FSHistoryService) SaveRequest(WeatherRequest) error {
	// TODO: Implement save request to file system
	// Append request to h.requests
	// Persist requests to file system (in memory for now, but could be reloaded from FS on restart of API)
	return nil
}

func (h *FSHistoryService) GetHistory(orderBy string, limit int) ([]WeatherRequest, error) {
	// TODO: Implement
	// Check h.Requests with provided ordering and return limit (0 for all)
	return nil, nil
}
