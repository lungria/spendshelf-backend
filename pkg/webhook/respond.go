package webhook

type errorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
