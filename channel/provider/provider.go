package provider

import (
	gp "github.com/krostar/nebulo-golib/provider"
	"github.com/krostar/nebulo-server/channel"
)

// Provider contains all the methods needed to manage channels
type Provider interface {
	gp.TablesManagement

	FindByID(ID int) (u *channel.Channel, err error)
	Update(u *channel.Channel, fields map[string]interface{}) (err error)
}

// P is the selected provider
var P Provider
