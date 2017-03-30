package provider

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/log"
)

type initDBFct func(dbmap *gorp.DbMap) (tableName string, err error)

// InitializeDatabase define the user table properties
func InitializeDatabase(db *sql.DB, dialect gorp.Dialect, config DefaultConfig, initDB initDBFct) (dbMap *gorp.DbMap, tableName string, err error) {
	dbMap = &gorp.DbMap{Db: db, Dialect: dialect}

	tableName, err = initDB(dbMap)
	if err != nil {
		return nil, tableName, fmt.Errorf("unable to initialize %q: %v", tableName, err)
	}

	if log.Verbosity == log.DEBUG {
		dbMap.TraceOn(strings.Title(tableName)+" Provider -", &ORPLogger{})
	}

	if config.DropTablesIfExists {
		err = dbMap.DropTablesIfExists()
		if err != nil {
			return nil, tableName, fmt.Errorf("unable to drop tables: %v", err)
		}
	}

	if config.CreateTablesIfNotExists {
		err = dbMap.CreateTablesIfNotExists()
		if err != nil {
			return nil, tableName, fmt.Errorf("unable to create tables: %v", err)
		}
	}
	return dbMap, tableName, nil
}
