package handler

import (
	"errors"

	"github.com/labstack/echo"
)

// AuthLogout is not handled
func AuthLogout(c echo.Context) error {
	panic(errors.New("not implemented"))
}
