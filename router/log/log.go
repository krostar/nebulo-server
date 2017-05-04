package log

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/krostar/nebulo-server/router/handler"

	"github.com/krostar/nebulo-golib/log"
	"github.com/labstack/echo"
)

// Request log a request with differents informations took from the request
func Request(c echo.Context, responseStatus int, duration time.Duration, responseSize int64) (err error) {
	var (
		req        = c.Request()
		loggedUser string
	)

	u, err := handler.GetLoggedUser(c.Get("user"))
	if err == nil {
		loggedUser, err = u.Repr()
		if err != nil {
			loggedUser = ""
		} else {
			loggedUser = fmt.Sprintf(" %s", loggedUser)
		}
	}

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

	uri := req.URL.RequestURI()
	uris := strings.Split(uri, "?")
	uris[0] = c.Path()
	uri = strings.Join(uris, "?")

	log.Requestf("%s -%s - \"%s %s\" %d %dms %s<>%d %q",
		remoteIP, loggedUser, req.Method, uri, responseStatus,
		duration.Nanoseconds()/1000000, rxBytes, responseSize, req.UserAgent())
	return nil
}
