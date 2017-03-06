package middleware

import (
	"errors"

	"github.com/krostar/nebulo/router/httperror"
	"github.com/labstack/echo"
)

var (
	errCertificateNotProvider = errors.New("authentication certificate not provided")
	errNoTLS                  = errors.New("authentication is based on TLS, without TLS authentication can't work")
)

// Auth handle the authentication process
func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if c.IsTLS() {
				if len(c.Request().TLS.PeerCertificates) == 1 {
					c.Set("userPublicKey", c.Request().TLS.PeerCertificates[0].PublicKey)
				} else {
					return httperror.HTTPUnauthorizedError(errCertificateNotProvider)
				}
			} else {
				httperror.HTTPInternalServerError(errNoTLS)
			}
			return next(c)
		}
	}
}
