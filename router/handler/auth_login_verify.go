package handler

import (
	"errors"

	"github.com/labstack/echo"
)

// AuthLoginVerify is not handled
func AuthLoginVerify(c echo.Context) (err error) {
	panic(errors.New("not implemented"))
}
