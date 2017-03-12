package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/krostar/nebulo/router/httperror"
	"github.com/krostar/nebulo/user"
	up "github.com/krostar/nebulo/user/provider"
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
					userCert := c.Request().TLS.PeerCertificates[0]

					c.Set("userCert", userCert)
					u, err := up.P.FindByPublicKey(userCert.PublicKeyAlgorithm, userCert.PublicKey)
					if err != nil {
						return httperror.HTTPUnauthorizedError(user.ErrNotFound)
					}

					now := time.Now()
					if u.LoginFirst.IsZero() {
						u.LoginFirst = now
					}
					u.LoginLast = now
					if err = up.P.Save(u); err != nil {
						return httperror.HTTPInternalServerError(fmt.Errorf("User save failed: %v", err))
					}
					c.Set("user", u)
				} else {
					return httperror.HTTPUnauthorizedError(errCertificateNotProvider)
				}
			} else {
				return httperror.HTTPInternalServerError(errNoTLS)
			}
			return next(c)
		}
	}
}
