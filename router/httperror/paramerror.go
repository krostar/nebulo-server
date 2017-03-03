package httperror

import "net/http"

// BadParam handle a bad request error
// caused by a bad parameter submission
func BadParam(reason interface{}) (herr *HTTPError) {
	return &HTTPError{
		Code:       http.StatusBadRequest,
		Err:        "parameter_bad",
		Parameters: reason,
	}
}
