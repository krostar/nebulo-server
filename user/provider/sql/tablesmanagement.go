package sql

import "github.com/krostar/nebulo-server/user"

// CreateTables create all the required tables for users
func (p *Provider) CreateTables() (err error) {
	u := &user.User{}
	ucm := &user.ChannelMembership{}

	return p.DB.CreateTable(ucm, u).Error
}

// DropTables delete all the users tables
func (p *Provider) DropTables() (err error) {
	u := &user.User{}
	ucm := &user.ChannelMembership{}

	return p.DB.DropTableIfExists(u, ucm).Error
}

// CreateIndexes create constrains and indexes on users tables
func (p *Provider) CreateIndexes() (err error) {
	ucm := &user.ChannelMembership{}

	channelMembershipModel := p.DB.Model(ucm)
	if err = channelMembershipModel.
		AddUniqueIndex("uniq_membership", "channel_id", "user_id").Error; err != nil {
		return err
	}
	return channelMembershipModel.
		AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").
		AddForeignKey("channel_id", "channels(id)", "CASCADE", "CASCADE").Error
}
