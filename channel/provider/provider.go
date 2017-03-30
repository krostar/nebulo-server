package provider

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/channel"
	"github.com/krostar/nebulo/log"
	gp "github.com/krostar/nebulo/provider"
)

// Provider represent the way we can interact with a provider
// to get informations about a channel
type Provider interface {
	SQLCreateQuery() (sqlCreationQuery string, err error)

	FindByID(ID int) (u *channel.Channel, err error)

	Update(u *channel.Channel, fields map[string]interface{}) (err error)
}

// P is the currently used provider
var P Provider

// Use set the new provider as the provider to use
func Use(newProvider Provider) (err error) {
	if P != nil {
		return errors.New("Hot database type change isn't supported")
	}
	P = newProvider

	return nil
}

// InitializeDatabase define the user table properties
func InitializeDatabase(db *sql.DB, dialect gorp.Dialect, dropTablesIfExists bool, createTablesIfNotExists bool) (dbmap *gorp.DbMap, channelTableName string, err error) {
	channelTableName = "channels"

	dbmap = &gorp.DbMap{Db: db, Dialect: dialect}

	channelTable := dbmap.AddTableWithName(channel.Channel{}, channelTableName)
	channelTable.SetKeys(true, "ID")

	if log.Verbosity == log.DEBUG {
		dbmap.TraceOn("Channel Provider -", &gp.ORPLogger{})
	}

	if dropTablesIfExists {
		err = dbmap.DropTablesIfExists()
		if err != nil {
			return nil, channelTableName, fmt.Errorf("unable to drop tables: %v", err)
		}
	}

	if createTablesIfNotExists {
		err = dbmap.CreateTablesIfNotExists()
		if err != nil {
			return nil, channelTableName, fmt.Errorf("unable to create tables: %v", err)
		}
	}
	return dbmap, channelTableName, nil
}
