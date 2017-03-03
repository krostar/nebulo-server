package handler

import (
	"errors"

	"github.com/labstack/echo"
)

// UserInfos is not handled
func UserInfos(c echo.Context) error {
	panic(errors.New("not implemented"))
}
