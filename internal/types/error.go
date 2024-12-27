package types

type InterServerError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
