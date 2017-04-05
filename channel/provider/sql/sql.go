package sql

import (
	"fmt"
	"reflect"

	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/channel"
	"github.com/krostar/nebulo/user/provider"
)

// Provider implement Interact and contain
// database useful variable
type Provider struct {
	provider.Provider
	DBMap            *gorp.DbMap
	ChannelTableName string
}

// SQLCreateQuery return the query used to generate the channels table
func (p *Provider) SQLCreateQuery() (string, error) {
	channelTable, err := p.DBMap.TableFor(reflect.TypeOf(channel.Channel{}), false)
	if err != nil {
		return "", fmt.Errorf("unable to get sql table for channel.Channel struct: %v", err)
	}
	return channelTable.SqlForCreate(true), nil
}

// Update only fiew fields from channel
func (p *Provider) Update(u *channel.Channel, fields map[string]interface{}) (err error) {
	// create SET <field>=<value> query from map
	var (
		sets     string
		args     []interface{}
		queryIdx = 0
	)
	for key, value := range fields {
		sets += " " + key + "=" + p.DBMap.Dialect.BindVar(queryIdx)
		args = append(args, value)
		queryIdx++
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%s", // nolint: gas
		p.ChannelTableName, sets, p.DBMap.Dialect.BindVar(queryIdx))
	if _, err = p.DBMap.Exec(query, append(args, u.ID)...); err != nil {
		return fmt.Errorf("unable to update channel informations: %v", err)
	}
	return nil
}
