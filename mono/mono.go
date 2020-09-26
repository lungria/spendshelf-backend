package mono

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient http.Client
}

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

type APIError struct {
	Description string `json:"errorDescription"`
}

func (err APIError) Error() string {
	return err.Description
}

// Time defines a timestamp encoded as epoch seconds in JSON
type Time time.Time

// MarshalJSON is used to convert the timestamp to JSON
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *Time) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return nil
}
