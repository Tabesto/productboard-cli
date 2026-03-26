package client

import "fmt"

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
		return &APIError{
			StatusCode: 401,
			Message:    "Authentication failed. Check your API token via `pboard configure` or PRODUCTBOARD_API_TOKEN env var.",
			ExitCode:   ExitAuthError,
		}
	case 403:
		return &APIError{
			StatusCode: 403,
			Message:    "Access denied. Your API token may not have permission for this resource.",
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
			msg = fmt.Sprintf("API error (%d): %s", statusCode, body)
		}
		return &APIError{
			StatusCode: statusCode,
			Message:    msg,
			ExitCode:   ExitGeneralError,
		}
	}
}

// NewNetworkError creates a network error.
func NewNetworkError(err error) *APIError {
	return &APIError{
		StatusCode: 0,
		Message:    fmt.Sprintf("Network error: %v. Check your internet connection.", err),
		ExitCode:   ExitGeneralError,
	}
}
