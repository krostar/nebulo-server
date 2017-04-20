package provider

import (
	gp "github.com/krostar/nebulo-golib/provider"
	"github.com/krostar/nebulo-server/channel"
	"github.com/krostar/nebulo-server/user"
)

// Provider contains all the methods needed to manage channels
type Provider interface {
	gp.TablesManagement

	Create(name string, creator user.User, members []user.User) (c *channel.Channel, err error)
	Find(toFind channel.Channel) (c *channel.Channel, err error)
	FindByID(ID int) (c *channel.Channel, err error)
}

// P is the selected provider
var P Provider
