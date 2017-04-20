package handler

import (
	"fmt"
	"net/http"

	"github.com/krostar/nebulo-golib/router/httperror"
	"github.com/labstack/echo"

	cp "github.com/krostar/nebulo-server/channel/provider"
	"github.com/krostar/nebulo-server/user"
	up "github.com/krostar/nebulo-server/user/provider"
)

// ChanCreateRequest store the request body for a ChanCreate request
type ChanCreateRequest struct {
	Name    string   `json:"name"`
	Members []string `json:"members_public_keys"`
}

// ChanCreate create a new channel for the logged user
func ChanCreate(c echo.Context) (err error) {
	u, err := GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	// bind the request body to the struct
	r := &ChanCreateRequest{}
	if err = c.Bind(r); err != nil {
		return httperror.HTTPBadRequestError(err)
	}

	// check all members existence and add them to a members array
	var (
		members []user.User
		member  *user.User
	)
	for _, m := range r.Members {
		member, err = up.P.FindByPublicKeyDERBase64(m)
		if err != nil {
			return httperror.New(http.StatusBadRequest, "members_public_keys",
				httperror.BadParam(fmt.Sprintf("unable to find member by public key: %v", err)),
			)
		}
		members = append(members, *member)
	}

	newChannel, err := cp.P.Create(r.Name, *u, members)
	if err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to create channel: %v", err))
	}

	return c.JSONPretty(http.StatusOK, newChannel, "    ")
}
