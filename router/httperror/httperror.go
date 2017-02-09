package httperror

import "fmt"
import "net/http"

// HTTPError contain the description of the error (type) and some useful values for this type
type HTTPError struct {
	Code       int         `json:"-"`
	Err        string      `json:"error"`
	Parameters interface{} `json:"parameters,omitempty"`
}

// ResponseHandler is the prototype of a ResponseHandler
type ResponseHandler func(interface{}) *HTTPError

// Error is useful to convert from HTTPError to string
func (e *HTTPError) Error() (err string) {
	err = "'" + e.Err + "'"
	if e.Parameters != nil {
		err += " with parameters '" + fmt.Sprint(e.Parameters) + "'"
	}
	return err
}

// HTTPUnsupportedMediaTypeError handle an HTTP error
func HTTPUnsupportedMediaTypeError(parameters interface{}) (herr *HTTPError) {
	if errType, ok := parameters.(error); ok {
		parameters = errType.Error()
	}
	return &HTTPError{
		Code:       http.StatusUnsupportedMediaType,
		Err:        "http_unsuported_media_type",
		Parameters: parameters,
	}
}

// HTTPBadRequestError handle an HTTP error
func HTTPBadRequestError(parameters interface{}) (herr *HTTPError) {
	if errType, ok := parameters.(error); ok {
		parameters = errType.Error()
	}
	return &HTTPError{
		Code:       http.StatusBadRequest,
		Err:        "http_bad_request",
		Parameters: parameters,
	}
}

// HTTPNotFoundError handle an HTTP error
func HTTPNotFoundError(parameters interface{}) (herr *HTTPError) {
	if errType, ok := parameters.(error); ok {
		parameters = errType.Error()
	}
	return &HTTPError{
		Code:       http.StatusNotFound,
		Err:        "http_not_found",
		Parameters: parameters,
	}
}

// HTTPUnauthorizedError handle an HTTP error
func HTTPUnauthorizedError(parameters interface{}) (herr *HTTPError) {
	if errType, ok := parameters.(error); ok {
		parameters = errType.Error()
	}
	return &HTTPError{
		Code:       http.StatusUnauthorized,
		Err:        "http_not_authorized",
		Parameters: parameters,
	}
}

// HTTPMethodNotAllowedError handle an HTTP error
func HTTPMethodNotAllowedError(parameters interface{}) (herr *HTTPError) {
	if errType, ok := parameters.(error); ok {
		parameters = errType.Error()
	}
	return &HTTPError{
		Code:       http.StatusMethodNotAllowed,
		Err:        "http_method_not_allowed",
		Parameters: parameters,
	}
}

// HTTPRequestEntityTooLargeError handle an HTTP error
func HTTPRequestEntityTooLargeError(parameters interface{}) (herr *HTTPError) {
	if errType, ok := parameters.(error); ok {
		parameters = errType.Error()
	}
	return &HTTPError{
		Code:       http.StatusRequestEntityTooLarge,
		Err:        "http_request_entity_too_large",
		Parameters: parameters,
	}
}

// HTTPInternalServerError handle an HTTP error
func HTTPInternalServerError(parameters interface{}) (herr *HTTPError) {
	if errType, ok := parameters.(error); ok {
		parameters = errType.Error()
	}
	return &HTTPError{
		Code:       http.StatusInternalServerError,
		Err:        "http_internal_server_error",
		Parameters: parameters,
	}
}
