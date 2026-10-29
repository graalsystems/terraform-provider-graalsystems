package graalsystems

import (
	"errors"
	"net/http"
	"strconv"

	sdk "github.com/graalsystems/sdk/go"
)

// isHTTPCodeError returns true if err is an http error with code statusCode
func isHTTPCodeError(err error, statusCode int) bool {
	if err == nil {
		return false
	}

	statusText := strconv.Itoa(statusCode) + " " + http.StatusText(statusCode)

	responseError := &sdk.GenericOpenAPIError{}
	if errors.As(err, &responseError) && responseError.Error() == statusText {
		return true
	}
	return false
}

// is404Error returns true if err is an HTTP 404 error
func is404Error(err error) bool {
	return isHTTPCodeError(err, http.StatusNotFound)
}

// toStringList converts a list of interfaces to a list of strings
func toStringList(input []interface{}) []string {
	var output []string
	for _, v := range input {
		if v == nil { // Maybe not a good practice, but useful for now
			v = ""
		}
		output = append(output, v.(string))
	}
	return output
}

// toStringMap converts a map of interfaces to a map of strings
func toStringMap(input map[string]interface{}) map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = v.(string)
	}
	return output
}

// deref returns the value of a string pointer, or "" when it is nil.
func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
