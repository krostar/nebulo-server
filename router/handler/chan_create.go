package handler

import (
	"crypto/sha256"
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
	Members []string `json:"members_public_key"`
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

	name, members, err := chanCreateChecks(*r)
	if err != nil {
		return err
	}

	newChannel, err := cp.P.Create(name, *u, members)
	if err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to create channel: %v", err))
	}

	return c.JSONPretty(http.StatusOK, newChannel, "    ")
}

func chanCreateChecks(ccr ChanCreateRequest) (name string, members []user.User, err error) {
	if len(ccr.Members) <= 0 {
		return "", nil, httperror.New(http.StatusBadRequest, "members_public_key",
			httperror.BadParam(fmt.Sprintf("no members in conversation")),
		)
	}

	// check all members existence and add them to a members array
	var (
		member      *user.User
		defaultName string
		repr        string
	)
	for _, m := range ccr.Members {
		member, err = up.P.FindByPublicKeyDERBase64(m)
		if err != nil {
			return "", nil, httperror.New(http.StatusBadRequest, "members_public_key",
				httperror.BadParam(fmt.Sprintf("unable to find member by public key: %v", err)),
			)
		}
		members = append(members, *member)
		repr, err = member.Repr()
		if err != nil {
			repr = member.FingerPrint
		}
		defaultName = fmt.Sprintf("%s - %s", defaultName, repr)
	}

	if ccr.Name == "" {
		hash := sha256.New()
		if _, err = hash.Write([]byte(defaultName)); err != nil {
			return "", nil, fmt.Errorf("unable to create sha256 hash of default name: %v", err)
		}
		name = fmt.Sprintf("%x", hash.Sum(nil))
	} else {
		name = ccr.Name
	}

	return name, members, err
}
