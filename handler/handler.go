package handler

import "errors"

var (
	// ErrUnableToSend is the message to write when the response can't be sent to client
	ErrUnableToSend = errors.New("unable to send response to client")
)
