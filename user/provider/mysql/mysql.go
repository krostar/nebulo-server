package sqlite

import (
	gp "github.com/krostar/nebulo/provider"
	"github.com/krostar/nebulo/user"
	"github.com/krostar/nebulo/user/provider"
	dp "github.com/krostar/nebulo/user/provider/sql"
)

// Provider implements the methods needed to manage a channel
// via a MySQL database
type Provider struct {
	dp.Provider
}

// Init initialize a MySQL provider and set it as the used provider
func Init() error {
	if gp.RP == nil {
		return gp.ErrRPIsNil
	}

	p := &Provider{}
	p.RootProvider = gp.RP

	provider.P = p
	return nil
}

// DropTables delete all the channels tables
func (p *Provider) DropTables() (err error) {
	u := &user.User{}
	ucm := &user.ChannelMembership{}

	err = gp.RP.DB.Exec("SET FOREIGN_KEY_CHECKS=0;").Error
	if err == nil {
		err = p.DB.DropTableIfExists(ucm, u).Error
	}
	if err == nil {
		err = gp.RP.DB.Exec("SET FOREIGN_KEY_CHECKS=1;").Error
	}
	return err
}
