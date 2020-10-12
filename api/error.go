package api

// Error describes standard API error response.
type Error struct {
	Message string `json:"errorMessage"`
}
