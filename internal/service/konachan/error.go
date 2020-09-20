package konachan

// https://konachan.com/help/api
type APIError struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}
