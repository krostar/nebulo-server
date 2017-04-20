package middleware

import (
	"time"

	"github.com/krostar/nebulo-server/router/log"
	"github.com/labstack/echo"
)

func mLog(next echo.HandlerFunc, c echo.Context) (err error) {
	res := c.Response()

	// get different useful information for logging purpose
	// 	execution time
	start := time.Now()
	if err = next(c); err != nil {
		c.Error(err)
		return err
	}
	stop := time.Now()

	return log.Request(c, res.Status, stop.Sub(start), res.Size)
}

// Log is the router middleware used to log request messages with the wanted format
func Log() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return mLog(next, c)
		}
	}
}
