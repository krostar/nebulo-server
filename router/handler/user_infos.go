package handler

import (
	"net/http"

	"github.com/krostar/nebulo/router/httperror"
	"github.com/labstack/echo"
)

// UserInfos is not handled
func UserInfos(c echo.Context) error {
	u, err := GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	return c.JSON(http.StatusOK, u)
}
