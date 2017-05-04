package sql

import "github.com/krostar/nebulo-server/message"

// CreateTables create all the required tables for messages
func (p *Provider) CreateTables() (err error) {
	m := &message.Message{}
	return p.DB.CreateTable(m).Error
}

// DropTables delete all the messages tables
func (p *Provider) DropTables() (err error) {
	m := &message.Message{}
	return p.DB.DropTableIfExists(m).Error
}

// CreateIndexes create constrains and indexes on messages tables
func (p *Provider) CreateIndexes() (err error) {
	m := &message.Message{}

	return p.DB.Model(m).
		AddForeignKey("channel_id", "users(id)", "CASCADE", "CASCADE").
		AddForeignKey("sender_id", "users(id)", "CASCADE", "CASCADE").
		AddForeignKey("receiver_id", "users(id)", "CASCADE", "CASCADE").Error
}
