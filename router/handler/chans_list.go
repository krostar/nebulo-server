package handler

import (
	"net/http"

	"github.com/krostar/nebulo-golib/router/httperror"
	"github.com/labstack/echo"

	cp "github.com/krostar/nebulo-server/channel/provider"
)

// ChansListRequest store the request body for a ChansList request
type ChansListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func ChansList(c echo.Context) (err error) {
	u, err := GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	// bind the request body to the struct
	r := &ChansListRequest{}
	if err = c.Bind(r); err != nil {
		return httperror.HTTPBadRequestError(err)
	}

	if r.Limit < 1 || r.Limit > 20 {
		r.Limit = 20
	}

	list, err := cp.P.List(*u, r.Offset, r.Limit)
	if err != nil {
		return httperror.HTTPInternalServerError(err)
	}

	return c.JSONPretty(http.StatusOK, list, "    ")
}
