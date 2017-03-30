package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	validator "gopkg.in/validator.v2"

	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/channel/provider"
	dp "github.com/krostar/nebulo/channel/provider/sql"
	gp "github.com/krostar/nebulo/provider"
	_ "github.com/mattn/go-sqlite3" // driver for database communication
)

// Provider inherit from default SQL Provider and contain
// configuration from a SQLite database
type Provider struct {
	dp.Provider
	config *gp.SQLiteConfig
}

// NewFromConfig return a new Provider based on the configuration
func NewFromConfig(config interface{}) (p *Provider, err error) {
	sqliteConfig, ok := config.(*gp.SQLiteConfig)
	if !ok {
		return nil, errors.New("unable to cast config to *sqlite.Config")
	}
	if err = validator.Validate(sqliteConfig); err != nil {
		return nil, fmt.Errorf("user file provider configuration validation failed: %v", err)
	}

	p = &Provider{
		config: sqliteConfig,
	}

	db, err := sql.Open("sqlite3", sqliteConfig.Filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to sqlite database: %v", err)
	}

	p.DBMap, p.ChannelTableName, err = provider.InitializeDatabase(db, &gorp.SqliteDialect{}, sqliteConfig.DropTablesIfExists, sqliteConfig.CreateTablesIfNotExists)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sqlite database: %v", err)
	}

	return p, nil
}
