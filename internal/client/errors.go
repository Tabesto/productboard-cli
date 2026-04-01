package client

import (
	"encoding/json"
	"fmt"
)

// Exit codes for CLI error types.
const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitAuthError    = 2
	ExitNotFound     = 3
	ExitRateLimited  = 4
	ExitInvalidInput = 5
)

// APIError represents an error from the ProductBoard API.
type APIError struct {
	StatusCode int
	Message    string
	ExitCode   int
}

func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates an appropriate error from an HTTP status code.
func NewAPIError(statusCode int, body string) *APIError {
	switch statusCode {
	case 401:
		msg := "Authentication failed. Check your API token via `pboard configure` or PRODUCTBOARD_API_TOKEN env var."
		if detail := extractV2ErrorDetail(body); detail != "" {
			msg = fmt.Sprintf("Authentication failed: %s", detail)
		}
		return &APIError{
			StatusCode: 401,
			Message:    msg,
			ExitCode:   ExitAuthError,
		}
	case 403:
		msg := "Access denied. Your API token may not have permission for this resource."
		if detail := extractV2ErrorDetail(body); detail != "" {
			msg = fmt.Sprintf("Access denied: %s", detail)
		}
		return &APIError{
			StatusCode: 403,
			Message:    msg,
			ExitCode:   ExitAuthError,
		}
	case 404:
		return &APIError{
			StatusCode: 404,
			Message:    "Resource not found.",
			ExitCode:   ExitNotFound,
		}
	case 429:
		return &APIError{
			StatusCode: 429,
			Message:    "Rate limit exceeded. Wait and retry, or reduce request frequency.",
			ExitCode:   ExitRateLimited,
		}
	default:
		if statusCode >= 500 {
			return &APIError{
				StatusCode: statusCode,
				Message:    fmt.Sprintf("ProductBoard API error (%d). Try again later.", statusCode),
				ExitCode:   ExitGeneralError,
			}
		}
		msg := fmt.Sprintf("API error (%d)", statusCode)
		if body != "" {
			// Try V2 error format first
			if detail := extractV2ErrorDetail(body); detail != "" {
				msg = fmt.Sprintf("API error (%d): %s", statusCode, detail)
			} else {
				msg = fmt.Sprintf("API error (%d): %s", statusCode, body)
			}
		}
		return &APIError{
			StatusCode: statusCode,
			Message:    msg,
			ExitCode:   ExitGeneralError,
		}
	}
}

// extractV2ErrorDetail attempts to parse a V2 error response and extract the detail message.
func extractV2ErrorDetail(body string) string {
	var v2Err struct {
		Errors []struct {
			Detail string `json:"detail"`
			Title  string `json:"title"`
		} `json:"errors"`
	}
	if err := json.Unmarshal([]byte(body), &v2Err); err != nil {
		return ""
	}
	if len(v2Err.Errors) == 0 {
		return ""
	}
	if v2Err.Errors[0].Detail != "" {
		return v2Err.Errors[0].Detail
	}
	return v2Err.Errors[0].Title
}

// NewNetworkError creates a network error.
func NewNetworkError(err error) *APIError {
	return &APIError{
		StatusCode: 0,
		Message:    fmt.Sprintf("Network error: %v. Check your internet connection.", err),
		ExitCode:   ExitGeneralError,
	}
}
