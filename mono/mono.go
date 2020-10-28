package mono

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client implements monobank API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient http.Client
}

// NewClient creates new instance of Client.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{baseURL: baseURL, apiKey: apiKey, httpClient: http.Client{}}
}

func (c *Client) performRequest(ctx context.Context, url, method string, requestBody []byte) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("cannot form request. %s", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Token", c.apiKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := response.Body.Close(); err == nil {
			err = fmt.Errorf("closing response body: %s", closeErr)
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %s", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		if len(body) == 0 {
			return nil, fmt.Errorf("request failed but no detailed error received. status code: %v", response.StatusCode)
		}

		apiErr := &APIError{}
		if err = json.Unmarshal(body, apiErr); err != nil {
			return nil, fmt.Errorf("failed unmarshal error form json body: %w", err)
		}

		return nil, apiErr
	}

	return body, nil
}

// APIError describes monobank API error.
type APIError struct {
	Description string `json:"errorDescription"`
}

// Error formats APIError to string.
func (err *APIError) Error() string {
	return err.Description
}
