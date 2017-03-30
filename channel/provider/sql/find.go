package sql

import (
	"database/sql"
	"fmt"

	"github.com/krostar/nebulo/channel"
	"github.com/krostar/nebulo/user"
)

// FindByID is used to find a channel from his ID
func (p *Provider) FindByID(id int) (u *channel.Channel, err error) {
	return p.find("id", id)
}

func (p *Provider) find(field string, value interface{}) (u *channel.Channel, err error) {
	u = new(channel.Channel)

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s=%s", // nolint: gas
		p.ChannelTableName, field, p.DBMap.Dialect.BindVar(0))
	if err = p.DBMap.SelectOne(u, query, value); err != nil {
		return nil, fmt.Errorf("unable to select channel in db: %v", err)
	}
	if err == sql.ErrNoRows || u == nil {
		return nil, user.ErrNotFound
	}

	return u, nil
}
