package sql

import "github.com/krostar/nebulo-server/channel"

// CreateTables create all the required tables for channels
func (p *Provider) CreateTables() (err error) {
	c := &channel.Channel{}

	if err = p.DB.CreateTable(c).Error; err != nil {
		return err
	}

	channelModel := p.DB.Model(c)
	return channelModel.Error
}

// DropTables delete all the channels tables
func (p *Provider) DropTables() (err error) {
	c := &channel.Channel{}

	return p.DB.DropTableIfExists(c).Error
}

// CreateIndexes create constrains and indexes on channels tables
func (p *Provider) CreateIndexes() (err error) {
	return nil
}
