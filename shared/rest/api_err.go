package rest

import "net/http"

// ApiError is an error type returned by the REST API.
// A message can be attached to a field (for form validation for example).
// Errors of this type will be returned as http.StatusBadRequest to the client by HandleErrors.
type ApiError struct {
	Message string `json:"message"`
	Field   string `json:"field"`
}

// NewApiError returns a new ApiError for given message and field.
func NewApiError(message, field string) *ApiError {
	return &ApiError{message, field}
}

// Error returns the error message.
func (err *ApiError) Error() string {
	return err.Message
}

// ErrorHandler is a special handler function to enable error handling.
type ErrorHandler func(http.ResponseWriter, *http.Request) []error

// ErrorMiddleware handles ErrorHandler errors.
func ErrorMiddleware(next ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		HandleErrors(w, r, err...)
	})
}
