package rest

import (
	"github.com/emvi/logbuch"
	"net/http"
	"syscall"
)

// WriteErrorResponse returns errors to the client with specified status code.
func WriteErrorResponse(w http.ResponseWriter, status int, errs ...error) {
	validationErrors := make(map[string]string, 0)
	apiErrors := make([]ApiError, 0)
	unexpectedErrors := make([]string, 0)

	for _, err := range errs {
		switch e := err.(type) {
		case *ApiError:
			if e.Field != "" {
				validationErrors[e.Field] = e.Message
			} else {
				apiErrors = append(apiErrors, *e)
			}
		case error:
			unexpectedErrors = append(unexpectedErrors, e.Error())
		default:
			logbuch.Error("Encountered error of unknown type", logbuch.Fields{"err": e})
		}
	}

	resp := struct {
		Validation map[string]string `json:"validation"`
		Errors     []ApiError        `json:"errors"`
		Exceptions []string          `json:"exceptions"`
	}{validationErrors, apiErrors, unexpectedErrors}

	w.WriteHeader(status)
	WriteResponse(w, resp)
}

// HandleErrors handles errors returned by REST endpoints.
// If the errors list contains a standard Go error, it will return status code 500 to the client.
// Else all errors passed must be of type ApiError or another custom type.
func HandleErrors(w http.ResponseWriter, r *http.Request, errs ...error) {
	// EPIPE ignore broken pipe errors
	if len(errs) == 1 && errs[0] == syscall.EPIPE {
		return
	}

	if len(errs) != 0 {
		logbuch.Debug("Errors on request", logbuch.Fields{"errs": errs, "url": r.URL.Path, "query": r.URL.RawQuery})

		if hasStandardError(errs) {
			WriteErrorResponse(w, http.StatusInternalServerError, errs...)
		} else {
			WriteErrorResponse(w, http.StatusBadRequest, errs...)
		}
	}
}

func hasStandardError(errs []error) bool {
	for _, err := range errs {
		switch err.(type) {
		case *ApiError:
			continue
		case error:
			return true
		}
	}

	return false
}
