package provider

import (
	"github.com/krostar/nebulo/channel"
	gp "github.com/krostar/nebulo/provider"
)

// Provider contains all the methods needed to manage channels
type Provider interface {
	gp.TablesManagement

	FindByID(ID int) (u *channel.Channel, err error)
	Update(u *channel.Channel, fields map[string]interface{}) (err error)
}

// P is the selected provider
var P Provider
