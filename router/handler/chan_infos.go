package handler

import (
	"net/http"

	"github.com/krostar/nebulo-golib/router/httperror"
	"github.com/labstack/echo"
)

func ChanInfos(c echo.Context) (err error) {
	_, err = GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	// TODO: we need same as chansList but for one channel

	return c.JSONPretty(http.StatusOK, nil, "    ")
}
