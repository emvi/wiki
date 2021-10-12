package rest

import (
	"encoding/json"
	"github.com/emvi/hide"
	iso6391 "github.com/emvi/iso-639-1"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	dateFormat      = "2006-01-02"
	defaultLangCode = "en"
)

var (
	supportedLangCodes = map[string]bool{"en": true, "de": true}
)

// DecodeJSON decodes a JSON object into the given object.
// If it fails to decode, it returns a new ApiError which can be send to the client.
func DecodeJSON(r *http.Request, req interface{}) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		logbuch.Debug("Error decoding JSON request", logbuch.Fields{"err": err})
		return NewApiError("Invalid format", "")
	}

	logbuch.Debug("Decoded JSON request", logbuch.Fields{"req": req})
	return nil
}

// WriteResponse converts and sends given object to the client as JSON.
// If it fails to encode, it responses with status code 500.
func WriteResponse(w http.ResponseWriter, resp interface{}) {
	// EPIPE ignore broken pipe errors
	if err := json.NewEncoder(w).Encode(resp); err != nil && err != syscall.EPIPE {
		logbuch.Error("Error marshalling response", logbuch.Fields{"err": err, "response": resp})
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetIntParam returns given request parameter as integer.
func GetIntParam(r *http.Request, name string) (int, error) {
	param := GetParam(r, name)

	if param == "" {
		return 0, nil
	}

	number, err := strconv.Atoi(param)

	if err != nil {
		logbuch.Debug("Error reading int parameter from url", logbuch.Fields{"err": err, "name": name})
		return 0, NewApiError("Invalid format", name)
	}

	return number, nil
}

// IdParam returns given request URL parameter as hide.ID.
// Example: /api/v1/test/{id} -> id
func IdParam(r *http.Request, name string) (hide.ID, error) {
	params := mux.Vars(r)
	id, err := hide.FromString(params[name])

	if err != nil {
		return 0, NewApiError("Invalid format", name)
	}

	return id, nil
}

// GetIntParam returns given request parameter as hide.ID.
func GetIdParam(r *http.Request, name string) (hide.ID, error) {
	id, err := hide.FromString(GetParam(r, name))

	if err != nil {
		return 0, NewApiError("Invalid format", name)
	}

	return id, nil
}

// GetIntParams returns given request parameter as []hide.ID.
// The IDs must be separated by ",". Example: ?param=1,2,3
func GetIdParams(r *http.Request, name string) ([]hide.ID, error) {
	param := GetParam(r, name)

	if param == "" {
		return []hide.ID{}, nil
	}

	str := strings.Split(param, ",")
	ids := make([]hide.ID, 0)

	for i := range str {
		id, err := hide.FromString(strings.TrimSpace(str[i]))

		if err != nil {
			return nil, NewApiError("Invalid format", name)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// GetIntParam returns given request parameter as bool.
func GetBoolParam(r *http.Request, name string) bool {
	param := strings.ToLower(GetParam(r, name))
	return param == "1" || param == "true"
}

// GetIntParam returns given request parameter as time.Time.
// The date format is specified as 2006-01-02.
func GetDateParam(r *http.Request, name string) (time.Time, error) {
	param := GetParam(r, name)

	if param == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse(dateFormat, param)

	if err != nil {
		logbuch.Debug("Error reading date parameter from url", logbuch.Fields{"err": err, "name": name})
		return time.Time{}, NewApiError("Invalid format", name)
	}

	return t, nil
}

// GetParams returns given request parameter as []string.
// The strings must be separated by ",". Example: ?param=one,two,three
func GetParams(r *http.Request, name string) []string {
	param := GetParam(r, name)

	if param == "" {
		return []string{}
	}

	str := strings.Split(param, ",")
	strs := make([]string, 0)

	for i := range str {
		strs = append(strs, str[i])
	}

	return strs
}

// GetIntParam returns given request parameter as string.
func GetParam(r *http.Request, name string) string {
	return strings.TrimSpace(r.URL.Query().Get(name))
}

// GetLangCode returns the language code in Accept-Language header.
func GetLangCode(r *http.Request) string {
	header := r.Header.Get("Accept-Language")
	parts := strings.Split(header, ";")

	if len(parts) == 0 || len(parts[0]) < 2 {
		return defaultLangCode
	}

	code := strings.ToLower(parts[0][:2])

	if iso6391.ValidCode(code) {
		return code
	}

	return defaultLangCode
}

// GetSupportedLangCode returns a language code from the Accept-Language header that is supported by the system.
func GetSupportedLangCode(r *http.Request) string {
	langCode := GetLangCode(r)
	_, ok := supportedLangCodes[langCode]

	if ok {
		return langCode
	}

	return defaultLangCode
}
