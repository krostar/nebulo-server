package middleware

import (
	"errors"
	"fmt"

	"github.com/krostar/nebulo-golib/router/httperror"
	"github.com/krostar/nebulo-golib/tools/cert"
	"github.com/labstack/echo"

	"github.com/krostar/nebulo-server/user"
	up "github.com/krostar/nebulo-server/user/provider"
)

var (
	errCertificateNotProvider = errors.New("authentication certificate not provided")
	errNoTLS                  = errors.New("authentication is based on TLS, without TLS authentication can't work")
)

func mAuth(next echo.HandlerFunc, c echo.Context) (err error) {
	// auth is based on certificate provided by clients during a TLS handshake
	if !c.IsTLS() {
		return httperror.HTTPInternalServerError(errNoTLS)
	}

	// if less or more than one certificate, its a bad request
	if len(c.Request().TLS.PeerCertificates) != 1 {
		return httperror.HTTPBadRequestError(errCertificateNotProvider)
	}

	userCert := c.Request().TLS.PeerCertificates[0]

	// check the certificate revokation
	revoked, err := cert.VerifyCertificate(userCert)
	if err != nil {
		return httperror.HTTPUnauthorizedError(fmt.Errorf("unable to verify certificate: %v", err))
	} else if revoked {
		return httperror.HTTPUnauthorizedError(errors.New("certificate is revoked"))
	}

	c.Set("userCert", userCert)
	u, err := up.P.FindByPublicKey(userCert.PublicKey)
	if err != nil {
		return httperror.HTTPUnauthorizedError(user.ErrNotFound)
	}

	if err = up.P.Login(u); err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("user save failed: %v", err))
	}
	c.Set("user", u)

	return next(c)
}

// Auth handle the authentication process
func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return mAuth(next, c)
		}
	}
}
