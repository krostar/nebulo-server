package handler

import (
	"errors"

	"github.com/labstack/echo"
)

// AuthLogin is not handled
func AuthLogin(c echo.Context) (err error) {
	panic(errors.New("not implemented"))
}
