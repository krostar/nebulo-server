package middleware

import (
	"net"
	"time"

	"github.com/krostar/nebulo/log"

	"github.com/labstack/echo"
)

// Log is the router middleware used to log request messages with the wanted format
func Log() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()

			// get different useful information for logging purpose
			// execution time
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			// active user
			// TODO: get active user repr from context
			user := "anon"

			// active user IP
			remoteIP := req.RemoteAddr
			if ip := req.Header.Get(echo.HeaderXRealIP); ip != "" {
				remoteIP = ip
			} else if ip = req.Header.Get(echo.HeaderXForwardedFor); ip != "" {
				remoteIP = ip
			} else {
				remoteIP, _, err = net.SplitHostPort(remoteIP)
				if err != nil {
					remoteIP = "unknown"
					log.Errorln("Unable to retrieve request client IP")
				}
			}

			// bytes read and writted
			rxBytes := req.Header.Get(echo.HeaderContentLength)
			if rxBytes == "" {
				rxBytes = "0"
			}

			log.Requestf("%s - %s - \"%s %s\" %d %dms %s<>%d %q %q", remoteIP,
				user, req.Method, req.URL.RequestURI(), res.Status,
				stop.Sub(start).Nanoseconds()/1000000, rxBytes, res.Size,
				req.Referer(), req.UserAgent())
			return
		}
	}
}
