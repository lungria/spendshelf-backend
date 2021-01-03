package api

import "net/http"

// Error describes standard API error response.
type Error struct {
	Message string `json:"errorMessage"`
}

// InternalServerError returns default error model for HTTP 500 responses.
func InternalServerError() (int, Error) {
	return http.StatusInternalServerError, Error{Message: "something went wrong"}
}
