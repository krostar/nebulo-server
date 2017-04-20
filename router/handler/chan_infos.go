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

	// TODO: we need, return last 20channels, offset to choose one
	/*
		- channel uniq repr
		- channel name
		- channel users participants
		- channel creation / last update
		- channel last message
	*/

	return c.JSONPretty(http.StatusOK, nil, "    ")
}
