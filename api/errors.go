package api

const (
	APIErrorCodeProviderUnavailable = "ERR_PROVIDER_UNAVAILABLE"
)

type APIError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type BadRequestError struct {
	Field       string `json:"field,omitempty"`
	Description string `json:"description,omitempty"`
}

func createAPIError(code, message string) APIError {
	return APIError{
		Code:    code,
		Message: message,
	}
}

func createBadRequestError(field, desc string) BadRequestError {
	return BadRequestError{
		Field:       field,
		Description: desc,
	}
}
