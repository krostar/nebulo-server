package httperror

import "net/http"

// UserExist handle a conflict error
// caused by an already existing user
func UserExist() (herr *HTTPError) {
	return &HTTPError{
		Code: http.StatusConflict,
		Err:  "user_exist",
	}
}

// UserNotFound handle a not found error
// caused by a non-existing user
func UserNotFound() (herr *HTTPError) {
	return &HTTPError{
		Code: http.StatusNotFound,
		Err:  "user_not_found",
	}
}
