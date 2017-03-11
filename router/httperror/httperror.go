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

// HTTPBadRequestError handle the 400 HTTP error
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

// HTTPUnauthorizedError handle the 401 HTTP error
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

// HTTPNotFoundError handle the 404 HTTP error
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

// HTTPMethodNotAllowedError handle the 405 HTTP error
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

// HTTPConflict handle the 409 HTTP error
func HTTPConflict(parameters interface{}) (herr *HTTPError) {
	if errType, ok := parameters.(error); ok {
		parameters = errType.Error()
	}
	return &HTTPError{
		Code:       http.StatusConflict,
		Err:        "http_conflict",
		Parameters: parameters,
	}
}

// HTTPRequestEntityTooLargeError handle the 413 HTTP error,
// also know as "Payload too large"
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

// HTTPUnsupportedMediaTypeError handle the 415 HTTP error
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

// HTTPInternalServerError handle the 500 HTTP error
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
