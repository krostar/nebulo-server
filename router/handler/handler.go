package handler

import (
	"errors"

	"github.com/krostar/nebulo/user"
)

var (
	// ErrUnableToSend is the message to write when the response can't be sent to client
	ErrUnableToSend = errors.New("unable to send response to client")
)

// GetLoggedUser return the current logged user based on the auth middleware
func GetLoggedUser(loggedUser interface{}) (u *user.User, err error) {
	u, ok := loggedUser.(*user.User)
	if !ok {
		return nil, errors.New("unable to cast user certificate to *user.User")
	}
	return u, nil
}
