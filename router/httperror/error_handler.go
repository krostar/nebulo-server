package httperror

import (
	"net/http"

	"github.com/krostar/nebulo-golib/router/httperror"

	"github.com/krostar/nebulo-golib/log"
	"github.com/labstack/echo"
)

// ErrorHandler handle all errors catched by echo allowing us to format the output
func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var errors *httperror.HTTPErrors

	switch herr := err.(type) {
	case *echo.HTTPError: // check for default echo http error
		// convert them to our custom errors
		convertFromEchoErrorToHEErrors := map[int]httperror.ResponseHandler{
			http.StatusUnsupportedMediaType:  httperror.HTTPUnsupportedMediaTypeError,
			http.StatusNotFound:              httperror.HTTPNotFoundError,
			http.StatusUnauthorized:          httperror.HTTPUnauthorizedError,
			http.StatusMethodNotAllowed:      httperror.HTTPMethodNotAllowedError,
			http.StatusRequestEntityTooLarge: httperror.HTTPRequestEntityTooLargeError,
			http.StatusBadRequest:            httperror.HTTPBadRequestError,
			http.StatusInternalServerError:   httperror.HTTPInternalServerError,
		}
		errors = httperror.New(herr.Code, "_", convertFromEchoErrorToHEErrors[herr.Code](nil))

	case *httperror.HTTPErrors: // then check for multiple errors
		// since it's our error format, we don't need to do anything
		errors = herr

	case *httperror.HTTPError: // then check for single error
		// create an error list from one error
		errors = httperror.New(herr.Code, "_", herr)

	default: // it's unexpected/unhandled
		if c.Echo().Debug {
			// if we are in local developpement, it's acceptable to return the unhandled error
			errors = httperror.New(http.StatusInternalServerError, "_", httperror.HTTPInternalServerError(err.Error()))
		} else {
			// in production, never return why we have this error, but log it
			log.Errorln("Unhandled internal error:", err)
			errors = httperror.New(http.StatusInternalServerError, "_", httperror.HTTPInternalServerError(nil))
		}
	}

	if err := c.JSON(errors.Code, errors); err != nil {
		panic(err)
	}
}
