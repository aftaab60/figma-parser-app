package errors

import (
	"encoding/json"
	"fmt"
)

type ErrorResponse struct {
	Err     string `json:"err"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HandleAPIError(statusCode int, body []byte) error {
	var errorResponse ErrorResponse
	// Try to parse the error response
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		return fmt.Errorf("HTTP %d: %s", statusCode, string(body))
	}

	switch statusCode {
	case 400:
		return fmt.Errorf("bad request: %s", errorResponse.Err)
	case 401:
		return fmt.Errorf("unauthorized: invalid Figma API token")
	case 403:
		return fmt.Errorf("forbidden: insufficient permissions for this file")
	case 404:
		return fmt.Errorf("not found: file does not exist or is not accessible")
	case 429:
		return fmt.Errorf("rate limit exceeded: too many requests")
	case 500:
		return fmt.Errorf("Figma server error: %s", errorResponse.Err)
	default:
		return fmt.Errorf("HTTP %d: %s", statusCode, errorResponse.Err)
	}
}
