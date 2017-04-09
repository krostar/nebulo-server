package provider

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/provider/orm"
)

// DefaultConfig is the configuration needed for running every database
type DefaultConfig struct {
	CreateTablesIfNotExists bool `json:"-"`
	DropTablesIfExists      bool `json:"-"`
}

// SQLiteConfig is the configuration needed for running an SQLite database
type SQLiteConfig struct {
	DefaultConfig
	Filepath string `json:"file" validate:"file=omitempty+readable:createifmissing"`
}

// MySQLConfig is the configuration needed for running an MySQL database
type MySQLConfig struct {
	DefaultConfig
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Database string `json:"database"`
}

// DatabaseType represent a database type
type DatabaseType int

const (
	// Unknown represent a non-handled database
	Unknown DatabaseType = iota
	// SQLite represent a SQLite database
	SQLite
	// MySQL represent a MySQL database
	MySQL
)

// TablesManagement contains all the required tables needed
// to manage tables
type TablesManagement interface {
	CreateTables() (err error)
	DropTables() (err error)
	CreateIndexes() (err error)
}

// RootProvider contains all the methods and properties needed
// to manage a database
type RootProvider struct {
	Type   DatabaseType
	DB     *gorm.DB
	Config interface{}
}

var (
	// RP is the selected database
	RP *RootProvider
	// ErrRPIsNil is the error throwed when RP is nil
	ErrRPIsNil = errors.New("database is nil")
)

// InitializeDatabase define the user table properties
func InitializeDatabase(p *RootProvider) {
	// in debug mode we want to see the sql requests
	if log.Verbosity == log.DEBUG {
		p.DB.SetLogger(&orm.Logger{})
		p.DB.LogMode(true)
	} else {
		p.DB.LogMode(false)
	}

	RP = p
}
