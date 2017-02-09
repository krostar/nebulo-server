package middleware

import (
	"fmt"
	"runtime"

	"github.com/krostar/nebulo/handler"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/router/httperror"

	"github.com/labstack/echo"
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// if we have a panic, dont crash and handle it properly
			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}

					// print stack
					stack := make([]byte, 64*1024)
					length := runtime.Stack(stack, true)
					log.Criticalln(err, "PANIC RECOVER", string(stack[:length]))

					// make a nebulo-http-compliant-error
					errr := httperror.HTTPInternalServerError(err)
					if err = c.JSON(errr.Code, httperror.New(errr.Code, "_", errr)); err != nil {
						log.Criticalln(handler.ErrUnableToSend, err)
					}
				}
			}()
			return next(c)
		}
	}
}
