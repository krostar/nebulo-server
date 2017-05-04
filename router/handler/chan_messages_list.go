package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/krostar/nebulo-golib/router/httperror"
	"github.com/labstack/echo"

	"github.com/krostar/nebulo-server/channel"
	cp "github.com/krostar/nebulo-server/channel/provider"
	mp "github.com/krostar/nebulo-server/message/provider"
)

func ChanMessagesList(c echo.Context) (err error) {
	u, err := GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	channelParam := c.Param("chan")

	queryParams := c.QueryParams()
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		return httperror.HTTPBadRequestError(fmt.Errorf("unable to parse limit %q: %v", queryParams.Get("limit"), err))
	}
	lastRead, err := time.Parse(time.RFC3339, queryParams.Get("last_read"))
	if err != nil {
		return httperror.HTTPBadRequestError(fmt.Errorf("unable to find last read %q: %v", queryParams.Get("last_read"), err))
	}

	if limit < -50 {
		limit = -50
	} else if limit > 50 {
		limit = 50
	}

	chann, err := cp.P.Find(channel.Channel{Name: channelParam})
	if err != nil {
		return httperror.HTTPBadRequestError(fmt.Errorf("unable to find channel %q: %v", channelParam, err))
	}
	list, err := mp.P.List(*u, *chann, lastRead, limit)
	if err != nil {
		return httperror.HTTPInternalServerError(err)
	}

	return c.JSONPretty(http.StatusOK, list, "    ")
}
