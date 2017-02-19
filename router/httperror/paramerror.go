package httperror

import "net/http"

// BadParam handle an HTTP error
func BadParam(reason interface{}) (herr *HTTPError) {
	return &HTTPError{
		Code:       http.StatusBadRequest,
		Err:        "parameter_bad",
		Parameters: reason,
	}
}
