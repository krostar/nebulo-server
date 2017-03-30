package provider

import (
	"errors"

	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/channel"
)

// Provider represent the way we can interact with a provider
// to get informations about a channel
type Provider interface {
	SQLCreateQuery() (sqlCreationQuery string, err error)

	FindByID(ID int) (u *channel.Channel, err error)

	Update(u *channel.Channel, fields map[string]interface{}) (err error)
}

// P is the currently used provider
var P Provider

// Use set the new provider as the provider to use
func Use(newProvider Provider) (err error) {
	if P != nil {
		return errors.New("Hot database type change isn't supported")
	}
	P = newProvider

	return nil
}

// InitializeDatabase define the channel table properties
func InitializeDatabase(dbmap *gorp.DbMap) (channelTableName string, err error) {
	channelTableName = "channels"

	channelTable := dbmap.AddTableWithName(channel.Channel{}, channelTableName)
	channelTable.SetKeys(true, "ID")

	return channelTableName, nil
}
