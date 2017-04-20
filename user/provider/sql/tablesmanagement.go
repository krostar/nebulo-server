package sql

import "github.com/krostar/nebulo-server/user"

// CreateTables create all the required tables for users
func (p *Provider) CreateTables() (err error) {
	u := &user.User{}
	return p.DB.CreateTable(u).Error
}

// DropTables delete all the users tables
func (p *Provider) DropTables() (err error) {
	u := &user.User{}
	return p.DB.DropTableIfExists(u).Error
}

// CreateIndexes create constrains and indexes on users tables
func (p *Provider) CreateIndexes() (err error) {
	u := &user.User{}

	userModel := p.DB.Model(u)
	return userModel.AddUniqueIndex("uniq_user", "key_public_der").Error
}
