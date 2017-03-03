package httperror

import (
	"net/http"

	"github.com/krostar/nebulo/log"
	"github.com/labstack/echo"
)

// ErrorHandler handle all errors catched by echo allowing us to format the output
func ErrorHandler(err error, c echo.Context) {
	if !c.Response().Committed {
		var errors *HTTPErrors

		if herr, ok := err.(*echo.HTTPError); ok { // check for default echo http error
			// convert them to our custom errors
			convertFromEchoErrorToHEErrors := map[int]ResponseHandler{
				http.StatusUnsupportedMediaType:  HTTPUnsupportedMediaTypeError,
				http.StatusNotFound:              HTTPNotFoundError,
				http.StatusUnauthorized:          HTTPUnauthorizedError,
				http.StatusMethodNotAllowed:      HTTPMethodNotAllowedError,
				http.StatusRequestEntityTooLarge: HTTPRequestEntityTooLargeError,
				http.StatusBadRequest:            HTTPBadRequestError,
				http.StatusInternalServerError:   HTTPInternalServerError,
			}
			errors = New(herr.Code, "_", convertFromEchoErrorToHEErrors[herr.Code](nil))
		} else if herr, ok := err.(*HTTPErrors); ok { // then check for multiple errors
			// since it's our error format, we don't need to do anything
			errors = herr
		} else if herr, ok := err.(*HTTPError); ok { // then check for single error
			// create an error list from one error
			errors = New(herr.Code, "_", herr)
		} else { // else it's unexpected/unhandled
			if c.Echo().Debug {
				// if we are in local developpement, it's acceptable to return the unhandled error
				errors = New(http.StatusInternalServerError, "_", HTTPInternalServerError(err.Error()))
			} else {
				// in production, never return why we have this error, but log it
				log.Errorln("Unhandled internal error:", err)
				errors = New(http.StatusInternalServerError, "_", HTTPInternalServerError(nil))
			}
		}

		if err := c.JSON(errors.Code, errors); err != nil {
			panic(err)
		}
	}
}
