package lightspeedsdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type APIError struct {
	StatusCode int
	Reason     string
	Detail     string
}

func parseAPIError(resp *http.Response) error {
	var apiErr APIError
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Reason:     "Failed to read response body",
		}
	}

	// Try to unmarshal the response body into the APIError struct.
	// This assumes the API might return structured error responses.
	// You may need to adjust this based on the actual error response structure of the API.
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Reason:     string(body),
		}
	}

	return &apiErr
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("API Error (%d): %s - %s", e.StatusCode, e.Reason, e.Detail)
	}
	return fmt.Sprintf("API Error (%d): %s", e.StatusCode, e.Reason)
}
