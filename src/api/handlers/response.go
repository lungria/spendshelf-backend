package handlers

type errorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ResponseFromError(err error) errorResponse {
	return errorResponse{
		Message: err.Error(),
		Error:   err.Error(),
	}
}

type messageResponse struct {
	Message string `json:"message"`
}
