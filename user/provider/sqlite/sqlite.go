package sqlite

import (
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"time"

	validator "gopkg.in/validator.v2"

	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/user"
	"github.com/krostar/nebulo/user/provider"
	_ "github.com/mattn/go-sqlite3" // driver for database communication
)

// Config is the configuration for the user provider file
type Config struct {
	Filepath                string `validate:"file=readable:createifmissing"`
	CreateTablesIfNotExists bool   `validate:"-"`
	DropTablesIfExists      bool   `validate:"-"`
}

// Provider implement provider.Provider and contain
// configuration from a SQLite database
type Provider struct {
	provider.Provider
	config *Config
	db     *sql.DB
	dbMap  *gorp.DbMap
}

// NewFromConfig return a new Provider based on the configuration
func NewFromConfig(sqliteConfig interface{}) (p *Provider, err error) {
	newConfig, ok := sqliteConfig.(*Config)
	if !ok {
		return nil, errors.New("unable to cast config to *sqlite.Config")
	}
	if err = validator.Validate(newConfig); err != nil {
		return nil, fmt.Errorf("user file provider configuration validation failed: %v", err)
	}

	db, dbmap, err := initializeDatabase(newConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sqlite database: %v", err)
	}

	return &Provider{
		config: newConfig,
		db:     db,
		dbMap:  dbmap,
	}, nil
}

func initializeDatabase(sqliteConfig *Config) (db *sql.DB, dbmap *gorp.DbMap, err error) {
	db, err = sql.Open("sqlite3", sqliteConfig.Filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to sqlite database: %v", err)
	}

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(user.User{}, "users").SetKeys(true, "ID")
	dbmap.TraceOn("User Provider - SQLITE -", provider.SQLLogger)

	if sqliteConfig.DropTablesIfExists {
		err = dbmap.DropTablesIfExists()
		if err != nil {
			return nil, nil, fmt.Errorf("unable to drop tables: %v", err)
		}
	}

	if sqliteConfig.CreateTablesIfNotExists {
		err = dbmap.CreateTablesIfNotExists()
		if err != nil {
			return nil, nil, fmt.Errorf("unable to create tables: %v", err)
		}
	}
	return db, dbmap, err
}

// Login update field on user login
func (p *Provider) Login(u *user.User) (err error) {
	if u == nil {
		return user.ErrUserNil
	}

	now := time.Now()

	u.LoginLast = now
	if u.LoginFirst.IsZero() {
		u.LoginFirst = now
	}

	if _, err = p.dbMap.Exec("UPDATE users SET login_first=?, login_last=? WHERE id=?", u.LoginFirst, u.LoginLast, u.ID); err != nil {
		return fmt.Errorf("unable to update login informations: %v", err)
	}
	return nil
}

// Register create a new User
func (p *Provider) Register(userToAdd *user.User) (u *user.User, err error) {
	// check user struct
	if userToAdd == nil {
		return nil, user.ErrUserNil
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

	err = p.dbMap.Insert(userToAdd)
	if err != nil {
		return nil, fmt.Errorf("unable to insert user: %v", err)
	}

	u = userToAdd
	return u, nil
}

// FindByPublicKey is used to find a user from his public key
func (p *Provider) FindByPublicKey(publicKeyAlgo x509.PublicKeyAlgorithm, publicKey interface{}) (u *user.User, err error) {
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal public key: %v", err)
	}

	u = new(user.User)
	if err = p.dbMap.SelectOne(u, "SELECT * FROM users WHERE key_public_der=?", publicKeyDER); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("unable to select user in db: %v", err)
	}
	if err == sql.ErrNoRows || u == nil {
		return nil, user.ErrNotFound
	}

	return u, nil
}

// FindByID is used to find a user from his ID
func (p *Provider) FindByID(ID int) (u *user.User, err error) {
	if err = p.dbMap.SelectOne(u, "SELECT * FROM users WHERE id=?", ID); err != nil {
		return nil, fmt.Errorf("unable to select user in db: %v", err)
	}
	if u == nil {
		return nil, user.ErrNotFound
	}

	return u, nil
}

// Save update every informations about a user
func (p *Provider) Save(u *user.User) (err error) {
	if _, err = p.dbMap.Update(u); err != nil {
		return fmt.Errorf("unable to save user: %v", err)
	}
	return nil
}
