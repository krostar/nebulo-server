package middleware

import (
	"errors"
	"fmt"

	"github.com/krostar/nebulo/router/httperror"
	"github.com/krostar/nebulo/tools/cert"
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

					if revoked, err := cert.VerifyCertificate(userCert); err != nil {
						return httperror.HTTPUnauthorizedError(fmt.Errorf("unable to verify certificate: %v", err))
					} else if revoked {
						return httperror.HTTPUnauthorizedError(errors.New("certificate is revoked"))
					}

					c.Set("userCert", userCert)
					u, err := up.P.FindByPublicKey(userCert.PublicKeyAlgorithm, userCert.PublicKey)
					if err != nil {
						return httperror.HTTPUnauthorizedError(user.ErrNotFound)
					}

					if err = up.P.Login(u); err != nil {
						return httperror.HTTPInternalServerError(fmt.Errorf("user save failed: %v", err))
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
