package provider

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	validator "gopkg.in/validator.v2"

	_ "github.com/go-sql-driver/mysql" // driver for database communication
	"github.com/krostar/nebulo/provider"
)

// Use initialize a MySQL provider and set it as the used provider
func Use(config interface{}) (err error) {
	mysqlConfig, ok := config.(*provider.MySQLConfig)
	if !ok {
		return errors.New("unable to cast config to *MySQLConfig")
	}
	if err = validator.Validate(mysqlConfig); err != nil {
		return fmt.Errorf("user file provider configuration validation failed: %v", err)
	}

	var credentials string
	if mysqlConfig.Password == "" {
		credentials = fmt.Sprintf("%s", mysqlConfig.Username)
	} else {
		credentials = fmt.Sprintf("%s:%s", mysqlConfig.Username, mysqlConfig.Password)
	}

	p := &provider.RootProvider{
		Config: mysqlConfig,
	}

	p.DB, err = gorm.Open("mysql", fmt.Sprintf("%s@%s/%s", credentials, mysqlConfig.Address, mysqlConfig.Database))
	if err != nil {
		return fmt.Errorf("unable to connect to mysql database: %v", err)
	}

	provider.InitializeDatabase(p)
	return nil
}
