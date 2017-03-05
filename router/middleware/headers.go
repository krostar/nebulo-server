package middleware

import (
	"github.com/krostar/nebulo/config"
	"github.com/labstack/echo"
)

// Headers return a middleware which add specific Headers
// to all responses
func Headers() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Config.TLSCertFile != "" {
				c.Response().Header().Add("Strict-Transport-Security", "max-age=63072000")
			}
			return next(c)
		}
	}
}
