package provider

import (
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/user"
)

// Provider represent the way we can interact with a provider
// to get informations about a user
type Provider interface {
	SQLCreateQuery() (sqlCreationQuery string, err error)

	Login(u *user.User) (err error)
	Create(userToAdd *user.User) (u *user.User, err error)
	Delete(u *user.User) (err error)

	FindByPublicKey(publicKeyAlgo x509.PublicKeyAlgorithm, publicKey interface{}) (u *user.User, err error)
	FindByID(ID int) (u *user.User, err error)

	Update(u *user.User, fields map[string]interface{}) (err error)
	Save(u *user.User) (err error)
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
func InitializeDatabase(db *sql.DB, dialect gorp.Dialect, dropTablesIfExists bool, createTablesIfNotExists bool) (dbmap *gorp.DbMap, err error) {
	dbmap = &gorp.DbMap{Db: db, Dialect: dialect}

	userTable := dbmap.AddTableWithName(user.User{}, "users")
	userTable.SetUniqueTogether("key_public_der", "key_public_algo")
	userTable.SetKeys(true, "ID")

	if log.Verbosity == log.DEBUG {
		dbmap.TraceOn("User Provider -", &ORPLogger{})
	}

	if dropTablesIfExists {
		err = dbmap.DropTablesIfExists()
		if err != nil {
			return nil, fmt.Errorf("unable to drop tables: %v", err)
		}
	}

	if createTablesIfNotExists {
		err = dbmap.CreateTablesIfNotExists()
		if err != nil {
			return nil, fmt.Errorf("unable to create tables: %v", err)
		}
	}
	return dbmap, nil
}
