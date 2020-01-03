package handlers

type errorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
