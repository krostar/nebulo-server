package handler

import (
	"fmt"
	"net/http"

	"github.com/krostar/nebulo-golib/router/httperror"
	"github.com/labstack/echo"

	"github.com/krostar/nebulo-server/channel"
	cp "github.com/krostar/nebulo-server/channel/provider"
	"github.com/krostar/nebulo-server/message"
	mp "github.com/krostar/nebulo-server/message/provider"
	"github.com/krostar/nebulo-server/user"
	up "github.com/krostar/nebulo-server/user/provider"
)

type messageInfos struct {
	Message  message.SecureMsg `json:"message"`
	Receiver string            `json:"receiver_pkey"`
}

// ChanMessageCreateRequest store the request body for a ChanMessageCreate request
type ChanMessageCreateRequest struct {
	ChannelName string         `json:"channel_name"`
	Messages    []messageInfos `json:"messages"`
}

func ChanMessageCreate(c echo.Context) (err error) {
	u, err := GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	// bind the request body to the struct
	r := ChanMessageCreateRequest{}
	if err = c.Bind(&r); err != nil {
		return httperror.HTTPBadRequestError(err)
	}

	var (
		receiver *user.User
		chnel    *channel.Channel
	)
	chnel, err = cp.P.FindByName(*u, r.ChannelName)
	if err != nil {
		return httperror.HTTPBadRequestError(fmt.Errorf("chan %q not found: %v", r.ChannelName, err))
	}

	for _, m := range r.Messages {
		receiver, err = up.P.FindByPublicKeyDERBase64(m.Receiver)
		if err != nil {
			return httperror.HTTPBadRequestError(fmt.Errorf("user not found: %v", err))
		}
		if _, err = mp.P.Create(*u, *receiver, *chnel, m.Message); err != nil {
			return httperror.HTTPInternalServerError(fmt.Errorf("unable to create channel: %v", err))
		}
	}

	return c.JSONPretty(http.StatusCreated, nil, "    ")
}
