package provider

import (
	"time"

	gp "github.com/krostar/nebulo-golib/provider"

	"github.com/krostar/nebulo-server/channel"
	"github.com/krostar/nebulo-server/message"
	"github.com/krostar/nebulo-server/user"
)

// Provider contains all the methods needed to manage channels
type Provider interface {
	gp.TablesManagement

	Create(sender user.User, receiver user.User, chann channel.Channel, msg message.SecureMsg) (m *message.Message, err error)
	List(receiver user.User, chann channel.Channel, lastRead time.Time, limit int) (m []*message.Message, err error)
}

// P is the selected provider
var P Provider
