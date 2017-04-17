package sql

import "github.com/krostar/nebulo-server/channel"

// CreateTables create all the required tables for channels
func (p *Provider) CreateTables() (err error) {
	c := &channel.Channel{}
	cum := &channel.UserMembership{}

	return p.DB.CreateTable(cum, c).Error
}

// DropTables delete all the channels tables
func (p *Provider) DropTables() (err error) {
	c := &channel.Channel{}
	cum := &channel.UserMembership{}

	return p.DB.DropTableIfExists(c, cum).Error
}

// CreateIndexes create constrains and indexes on channels tables
func (p *Provider) CreateIndexes() (err error) {
	cum := &channel.UserMembership{}

	userMembershipModel := p.DB.Model(cum)
	if err = userMembershipModel.
		AddUniqueIndex("uniq_membership", "channel_id", "user_id").Error; err != nil {
		return err
	}
	return userMembershipModel.
		AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").
		AddForeignKey("channel_id", "channels(id)", "CASCADE", "CASCADE").Error
}
