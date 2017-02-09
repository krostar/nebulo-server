package httperror

import "errors"

// HTTPErrors contain all the error list
type HTTPErrors struct {
	Code   int                   `json:"-"`
	Errors map[string]*HTTPError `json:"errors"`
}

var (
	errHEnil = errors.New("HTTPErrors is nil")
)

// New creare errors and add the item error
func New(code int, typ string, err *HTTPError) (errors *HTTPErrors) {
	e := new(HTTPErrors)
	e.Errors = make(map[string]*HTTPError)

	e.Code = code
	e = e.Add(typ, err)

	return e
}

// Add add in errors the item error
func (e *HTTPErrors) Add(typ string, err *HTTPError) (errors *HTTPErrors) {
	if e.Errors == nil {
		panic(errHEnil)
	}

	e.Errors[typ] = err
	return e
}

// Error concat all the errors in one error
func (e *HTTPErrors) Error() (err string) {
	if e.Errors == nil {
		panic(errHEnil)
	}

	err = ""
	for key := range e.Errors {
		err += key + ":" + e.Errors[key].Error() + "\n"
	}
	return err
}
