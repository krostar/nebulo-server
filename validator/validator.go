package validator

import (
	"net/http"
	"strings"

	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/router/httperror"

	validator "gopkg.in/validator.v2"
)

func init() {
	// tell to the validator lib which function to use for our custom validators
	var err error
	if err = validator.SetValidationFunc("file", File); err != nil {
		log.Criticalln(err)
		panic(err)
	}
	if err = validator.SetValidationFunc("nonempty", NonEmpty); err != nil {
		log.Criticalln(err)
		panic(err)
	}
}

// HTTPErrors transform lib-errors into send-able HTTPErrors
func HTTPErrors(err error) *httperror.HTTPErrors {
	// get the lib-error-map
	vErrs, ok := err.(validator.ErrorMap)
	if !ok {
		return nil
	}

	// transform the map into a list of httperror
	var httpErrors *httperror.HTTPErrors
	for key, errs := range vErrs {
		var erra []string
		for _, err := range errs {
			erra = append(erra, err.Error())
		}
		key = strings.ToLower(key)
		if httpErrors == nil {
			httpErrors = httperror.New(http.StatusBadRequest, key, httperror.BadParam(erra))
		} else {
			httpErrors = httpErrors.Add(key, httperror.BadParam(erra))
		}
	}

	return httpErrors
}
