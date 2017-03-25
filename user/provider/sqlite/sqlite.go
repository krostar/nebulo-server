package sqlite

import (
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	validator "gopkg.in/validator.v2"

	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/user"
	"github.com/krostar/nebulo/user/provider"
	_ "github.com/mattn/go-sqlite3" // we use this sql driver for database communication
)

// Config is the configuration for the user provider file
type Config struct {
	Filepath                string `validate:"file=readable:createifmissing"`
	CreateTablesIfNotExists bool   `validate:"-"`
	DropTablesIfExists      bool   `validate:"-"`
}

// ProviderSQLite implement Provider and contain the configuration
// and the user list loaded from a SQLite database
type ProviderSQLite struct {
	provider.Provider
	config *Config
	db     *sql.DB
	dbmap  *gorp.DbMap
}

// NewFromConfig return a new ProviderSQLite based on the configuration
func NewFromConfig(config interface{}) (p *ProviderSQLite, err error) {
	newConfig, ok := config.(*Config)
	if !ok {
		return nil, errors.New("unable to cast config to *sqlite.Config")
	}
	if err = validator.Validate(config); err != nil {
		return nil, fmt.Errorf("user file provider configuration validation failed: %v", err)
	}

	db, dbmap, err := initializeDatabase(newConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sqlite database: %v", err)
	}

	return &ProviderSQLite{
		config: newConfig,
		db:     db,
		dbmap:  dbmap,
	}, nil
}

func initializeDatabase(config *Config) (db *sql.DB, dbmap *gorp.DbMap, err error) {
	db, err = sql.Open("sqlite3", config.Filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to sqlite database: %v", err)
	}

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(user.User{}, "users").SetKeys(true, "ID")
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))

	if config.DropTablesIfExists {
		err = dbmap.DropTablesIfExists()
		if err != nil {
			return nil, nil, fmt.Errorf("unable to drop tables: %v", err)
		}
	}

	if config.CreateTablesIfNotExists {
		err = dbmap.CreateTablesIfNotExists()
		if err != nil {
			return nil, nil, fmt.Errorf("unable to create tables: %v", err)
		}
	}
	return db, dbmap, err
}

// Register create a new User
func (p *ProviderSQLite) Register(userToAdd *user.User) (u *user.User, err error) {
	if p == nil {
		return nil, errors.New("sqlite provider is nil")
	}

	// check user struct
	if userToAdd == nil {
		return nil, errors.New("user is nil")
	}

	// check if user exist
	userPublicKey, err := x509.ParsePKIXPublicKey(userToAdd.PublicKeyDER)
	if err != nil {
		return nil, fmt.Errorf("unable to parse user public key: %v", err)
	}
	_, err = p.FindByPublicKey(userToAdd.PublicKeyAlgorithm, userPublicKey)
	if err != nil && err != user.ErrNotFound {
		return nil, fmt.Errorf("unable to find user: %v", err)
	}
	if err != user.ErrNotFound {
		return nil, errors.New("an user already exist with this public key")
	}

	userToAdd.SignUp = time.Now()

	err = p.dbmap.Insert(userToAdd)
	if err != nil {
		return nil, fmt.Errorf("unable to insert user: %v", err)
	}

	u = userToAdd
	return u, nil
}

// FindByPublicKey is used to find a user from his public key
func (p *ProviderSQLite) FindByPublicKey(publicKeyAlgo x509.PublicKeyAlgorithm, publicKey interface{}) (u *user.User, err error) {
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal public key: %v", err)
	}

	u = new(user.User)
	if err = p.dbmap.SelectOne(u, "SELECT * FROM users WHERE key_public_der=?", publicKeyDER); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("unable to select user in db: %v", err)
	}
	if err == sql.ErrNoRows || u == nil {
		return nil, user.ErrNotFound
	}

	return u, nil
}

// FindByID is used to find a user from his ID
func (p *ProviderSQLite) FindByID(ID int) (u *user.User, err error) {
	if err = p.dbmap.SelectOne(u, "SELECT * FROM users WHERE id=?", ID); err != nil {
		return nil, fmt.Errorf("unable to select user in db: %v", err)
	}
	if u == nil {
		return nil, user.ErrNotFound
	}

	return u, nil
}

// Save modifications made on user
func (p *ProviderSQLite) Save(u *user.User) (err error) {
	// TODO: maybe a bit hard, nope ?
	_, err = p.dbmap.Update(u)
	if err != nil {
		return fmt.Errorf("unable to save user: %v", err)
	}
	return nil
}
