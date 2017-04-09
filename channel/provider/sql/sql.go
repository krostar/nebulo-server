package sql

import (
	"fmt"

	"github.com/krostar/nebulo/channel"
	"github.com/krostar/nebulo/channel/provider"
	gp "github.com/krostar/nebulo/provider"
)

// Provider implements the methods needed to manage channels
// for every SQL based database
type Provider struct {
	*gp.RootProvider
	provider.Provider
}

// Update update only fiew fields from user
func (p *Provider) Update(c *channel.Channel, fields map[string]interface{}) (err error) {
	if c == nil {
		return channel.ErrNil
	}

	if err = p.DB.Model(c).Updates(fields).Error; err != nil {
		return fmt.Errorf("unable to update channel informations: %v", err)
	}
	return nil
}
