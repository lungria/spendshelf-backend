package handlers

type errorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func responseFromError(err error, message string) errorResponse {
	return errorResponse{
		Message: message,
		Error:   err.Error(),
	}
}

type messageResponse struct {
	Message string `json:"message"`
}
