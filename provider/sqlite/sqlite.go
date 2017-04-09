package sqlite

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	validator "gopkg.in/validator.v2"

	"github.com/krostar/nebulo/provider"
	_ "github.com/mattn/go-sqlite3" // driver for database communication
)

// Use initialize a SQLite provider and set it as the used provider
func Use(config interface{}) (err error) {
	sqliteConfig, ok := config.(*provider.SQLiteConfig)
	if !ok {
		return errors.New("unable to cast config to *SQLiteConfig")
	}
	if err = validator.Validate(sqliteConfig); err != nil {
		return fmt.Errorf("sqlite provider configuration validation failed: %v", err)
	}

	p := &provider.RootProvider{
		Config: sqliteConfig,
	}

	p.DB, err = gorm.Open("sqlite3", sqliteConfig.Filepath)
	if err != nil {
		return fmt.Errorf("unable to connect to sqlite database: %v", err)
	}

	provider.InitializeDatabase(p)
	return nil
}
