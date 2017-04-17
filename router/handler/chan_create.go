package handler

import (
	"net/http"

	"github.com/krostar/nebulo-golib/router/httperror"
	"github.com/labstack/echo"
)

func ChanCreate(c echo.Context) (err error) {
	_, err = GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	// TODO: we need:
	/*
		- channel name (optionnel)
		- channel user(s) participant(s)
	*/

	return c.JSONPretty(http.StatusOK, nil, "    ")
}
