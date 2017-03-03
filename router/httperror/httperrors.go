package httperror

import "errors"

// HTTPErrors contain all the errors
type HTTPErrors struct {
	Code   int                   `json:"-"`
	Errors map[string]*HTTPError `json:"errors"`
}

var (
	// ErrHEIsNil is an error returned
	// when the claims object is nil
	ErrHEIsNil = errors.New("HTTPErrors is nil")
)

// New create a new error list and add the item error
func New(code int, typ string, err *HTTPError) (errors *HTTPErrors) {
	e := new(HTTPErrors)
	e.Errors = make(map[string]*HTTPError)

	e.Code = code
	e = e.Add(typ, err)

	return e
}

// Add add in error list the error
func (e *HTTPErrors) Add(typ string, err *HTTPError) (errors *HTTPErrors) {
	if e.Errors == nil {
		panic(ErrHEIsNil)
	}

	e.Errors[typ] = err
	return e
}

// Error concat all the errors in one message
func (e *HTTPErrors) Error() (err string) {
	if e.Errors == nil {
		panic(ErrHEIsNil)
	}

	err = ""
	for key := range e.Errors {
		err += key + ":" + e.Errors[key].Error() + "\n"
	}
	return err
}
