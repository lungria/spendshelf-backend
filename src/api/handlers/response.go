package handlers

type errorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type messageResponse struct {
	Message string `json:"message"`
}
