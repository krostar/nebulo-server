package middleware

import (
	"github.com/labstack/echo"
)

func mMisc(next echo.HandlerFunc, c echo.Context) (err error) {
	if c.Request().TLS != nil {
		c.Response().Header().Add("Strict-Transport-Security", "max-age=63072000")
	}
	return next(c)
}

// Misc return a middleware which add things to all responses
func Misc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return mMisc(next, c)
		}
	}
}
