package sqlite

import (
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
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
	config    *Config
	dbMap     *gorp.DbMap
	tableName string
}

// NewFromConfig return a new Provider based on the configuration
func NewFromConfig(config interface{}) (p *Provider, err error) {
	sqliteConfig, ok := config.(*Config)
	if !ok {
		return nil, errors.New("unable to cast config to *sqlite.Config")
	}
	if err = validator.Validate(sqliteConfig); err != nil {
		return nil, fmt.Errorf("user file provider configuration validation failed: %v", err)
	}

	db, err := sql.Open("sqlite3", sqliteConfig.Filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to sqlite database: %v", err)
	}

	dbmap, err := provider.InitializeDatabase(db, &gorp.SqliteDialect{}, sqliteConfig.DropTablesIfExists, sqliteConfig.CreateTablesIfNotExists)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sqlite database: %v", err)
	}

	userTable, err := dbmap.TableFor(reflect.TypeOf(user.User{}), false)
	if err != nil {
		return nil, fmt.Errorf("unable to get user table name: %v", err)
	}

	return &Provider{
		config:    sqliteConfig,
		dbMap:     dbmap,
		tableName: userTable.TableName,
	}, nil
}

// SQLCreateQuery return the query used to generate the users table
func (p *Provider) SQLCreateQuery() (string, error) {
	userTable, err := p.dbMap.TableFor(reflect.TypeOf(user.User{}), false)
	if err != nil {
		return "", fmt.Errorf("unable to get sql table for user.User struct: %v", err)
	}
	return userTable.SqlForCreate(true), nil
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

	if _, err = p.dbMap.Exec("UPDATE "+p.tableName+" SET login_first=?, login_last=? WHERE id=?", u.LoginFirst, u.LoginLast, u.ID); err != nil {
		return fmt.Errorf("unable to update login informations: %v", err)
	}

	return nil
}

// Create a new user
func (p *Provider) Create(userToAdd *user.User) (u *user.User, err error) {
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

// Delete a existing user
func (p *Provider) Delete(u *user.User) (err error) {
	if u == nil {
		return user.ErrUserNil
	}

	// check if user exist
	_, err = p.FindByID(u.ID)
	if err != nil && err != user.ErrNotFound {
		return fmt.Errorf("unable to find user: %v", err)
	} else if err == user.ErrNotFound {
		return user.ErrNotFound
	}

	_, err = p.dbMap.Delete(u)
	if err != nil {
		return fmt.Errorf("unable to delete user: %v", err)
	}

	return nil
}

// FindByPublicKey is used to find a user from his public key
func (p *Provider) FindByPublicKey(publicKeyAlgo x509.PublicKeyAlgorithm, publicKey interface{}) (u *user.User, err error) {
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal public key: %v", err)
	}

	u = new(user.User)
	if err = p.dbMap.SelectOne(u, "SELECT * FROM "+p.tableName+" WHERE key_public_der=?", publicKeyDER); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("unable to select user in db: %v", err)
	}
	if err == sql.ErrNoRows || u == nil {
		return nil, user.ErrNotFound
	}

	return u, nil
}

// FindByID is used to find a user from his ID
func (p *Provider) FindByID(ID int) (u *user.User, err error) {
	u = new(user.User)
	if err = p.dbMap.SelectOne(u, "SELECT * FROM "+p.tableName+" WHERE id=?", ID); err != nil {
		return nil, fmt.Errorf("unable to select user in db: %v", err)
	}
	if u == nil {
		return nil, user.ErrNotFound
	}

	return u, nil
}

func (p *Provider) Update(u *user.User, fields map[string]interface{}) (err error) {
	var sets string
	var args []interface{}
	for key, value := range fields {
		sets += " " + key + "=?"
		args = append([]interface{}{value}, args...)
	}

	if _, err = p.dbMap.Exec("UPDATE "+p.tableName+" SET"+sets, args...); err != nil {
		return fmt.Errorf("unable to update user informations: %v", err)
	}
	return nil
}

// Save update every informations about a user
func (p *Provider) Save(u *user.User) (err error) {
	if _, err = p.dbMap.Update(u); err != nil {
		return fmt.Errorf("unable to save user: %v", err)
	}
	return nil
}
