package validator

import (
	"errors"
	"reflect"

	validator "gopkg.in/validator.v2"
)

// NonEmpty look for a non empty string for the validator tag
func NonEmpty(v interface{}, param string) (err error) {
	// get the value to test
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validator.ErrUnsupported
	}

	if len(st.String()) == 0 {
		return errors.New("missing or empty")
	}

	return nil
}
