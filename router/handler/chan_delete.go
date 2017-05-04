package handler

import (
	"net/http"

	"github.com/krostar/nebulo-golib/router/httperror"
	"github.com/labstack/echo"
)

func ChanDelete(c echo.Context) (err error) {
	_, err = GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	return c.JSONPretty(http.StatusOK, nil, "    ")
}
