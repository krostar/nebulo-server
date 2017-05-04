package sqlite

import (
	gp "github.com/krostar/nebulo-golib/provider"
	"github.com/krostar/nebulo-server/message/provider"
	dp "github.com/krostar/nebulo-server/message/provider/sql"
)

// Provider implements the methods needed to manage a channel
// via a SQLite database
type Provider struct {
	dp.Provider
}

// Init initialize a SQLite provider and set it as the used provider
func Init() error {
	if gp.RP == nil {
		return gp.ErrRPIsNil
	}

	p := &Provider{}
	p.RootProvider = gp.RP

	provider.P = p
	return nil
}

// CreateIndexes create constrains and indexes on channels tables
func (p *Provider) CreateIndexes() (err error) {
	return nil
}
