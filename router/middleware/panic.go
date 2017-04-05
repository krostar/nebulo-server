package middleware

import (
	"encoding/base64"
	"fmt"
	"runtime"

	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/router/handler"
	"github.com/krostar/nebulo/router/httperror"

	"github.com/labstack/echo"
)

func mRecover(next echo.HandlerFunc, c echo.Context) (err error) {
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

			// error stack
			stack := make([]byte, 64*1024)
			length := runtime.Stack(stack, true)

			var printableStack string
			var errHE *httperror.HTTPError

			// useful informations can be shown for debug pupose
			if c.Echo().Debug {
				printableStack = "\n" + string(stack[:length])
				errHE = httperror.HTTPInternalServerError(fmt.Errorf("panic recover: %v", err))
			} else {
				printableStack = base64.StdEncoding.EncodeToString(stack[:length])
				errHE = httperror.HTTPInternalServerError(nil)
			}

			log.Logln(log.CRITICAL, 5, err, printableStack)

			// return a nebulo http compliant error
			if err = c.JSON(errHE.Code, httperror.New(errHE.Code, "_", errHE)); err != nil {
				log.Criticalln(handler.ErrUnableToSend, err)
			}
		}
	}()
	return next(c)
}

// Recover return a middleware which recover from panics
// anywhere in the chain and handle the error
// with the centralized custom errors handler.
func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return mRecover(next, c)
		}
	}
}
