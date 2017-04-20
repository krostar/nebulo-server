package sql

import (
	"fmt"

	"github.com/krostar/nebulo-server/channel"
)

// Find is used to find a channel from the setted field
func (p *Provider) Find(toFind channel.Channel) (c *channel.Channel, err error) {
	c = new(channel.Channel)

	if p.DB.Where(toFind).First(c).RecordNotFound() {
		return nil, channel.ErrNotFound
	}
	if err = p.DB.Error; err != nil {
		return nil, fmt.Errorf("unable to select channel in db: %v", err)
	}

	return c, nil
}

// FindByID is used to find a channel from his ID
func (p *Provider) FindByID(id int) (c *channel.Channel, err error) {
	return p.find("id", id)
}

func (p *Provider) find(field string, value interface{}) (c *channel.Channel, err error) {
	c = new(channel.Channel)

	if p.DB.Where(field+" = ?", value).First(c).RecordNotFound() {
		return nil, channel.ErrNotFound
	}
	if err = p.DB.Error; err != nil {
		return nil, fmt.Errorf("unable to select channel in db: %v", err)
	}

	return c, nil
}
