package handlers

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ResponseFromError(err error) ErrorResponse {
	return ErrorResponse{
		Message: err.Error(),
		Error:   err.Error(),
	}
}

type messageResponse struct {
	Message string `json:"message"`
}
