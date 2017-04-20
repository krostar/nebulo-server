package sql

import (
	"fmt"

	gp "github.com/krostar/nebulo-golib/provider"
	"github.com/krostar/nebulo-server/channel"
	"github.com/krostar/nebulo-server/channel/provider"
	"github.com/krostar/nebulo-server/user"
)

// Provider implements the methods needed to manage channels
// for every SQL based database
type Provider struct {
	*gp.RootProvider
	provider.Provider
}

// Create create a channel if needed, or return an exsting one with the same requirements
func (p *Provider) Create(name string, creator user.User, members []user.User) (c *channel.Channel, err error) {
	c = &channel.Channel{
		Name:    name,
		Creator: creator,
	}

	oldChannel, err := p.Find(*c)
	if err != nil && err != channel.ErrNotFound {
		return nil, fmt.Errorf("unable to find channel: %v", err)
	}
	if err != nil {
		return oldChannel, nil
	}

	if err = p.DB.Create(c).Error; err != nil {
		return nil, fmt.Errorf("unable to insert channel: %v", err)
	}

	return c, nil
}
