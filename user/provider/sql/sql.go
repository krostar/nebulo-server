package sql

import (
	"crypto/x509"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/user"
	"github.com/krostar/nebulo/user/provider"
)

// Provider implement Interact and contain
// database useful variable
type Provider struct {
	provider.Provider
	DBMap         *gorp.DbMap
	UserTableName string
}

// SQLCreateQuery return the query used to generate the users table
func (p *Provider) SQLCreateQuery() (string, error) {
	userTable, err := p.DBMap.TableFor(reflect.TypeOf(user.User{}), false)
	if err != nil {
		return "", fmt.Errorf("unable to get sql table for user.User struct: %v", err)
	}
	return userTable.SqlForCreate(true), nil
}

// Login update field on user login
func (p *Provider) Login(u *user.User) (err error) {
	if u == nil {
		return user.ErrNil
	}

	now := time.Now()

	u.LoginLast = now
	if u.LoginFirst.IsZero() {
		u.LoginFirst = now
	}

	query := fmt.Sprintf("UPDATE %s SET login_first=%s, login_last=%s WHERE id=%s", // nolint: gas
		p.UserTableName, p.DBMap.Dialect.BindVar(0), p.DBMap.Dialect.BindVar(1), p.DBMap.Dialect.BindVar(2))
	if _, err = p.DBMap.Exec(query, u.LoginFirst, u.LoginLast, u.ID); err != nil {
		return fmt.Errorf("unable to update login informations: %v", err)
	}

	return nil
}

// Create a new user
func (p *Provider) Create(userToAdd *user.User) (u *user.User, err error) {
	if userToAdd == nil {
		return nil, user.ErrNil
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

	err = p.DBMap.Insert(userToAdd)
	if err != nil {
		return nil, fmt.Errorf("unable to insert user: %v", err)
	}

	u = userToAdd
	return u, nil
}

// Delete a existing user
func (p *Provider) Delete(u *user.User) (err error) {
	if u == nil {
		return user.ErrNil
	}

	// check if user exist
	_, err = p.FindByID(u.ID)
	if err != nil && err != user.ErrNotFound {
		return fmt.Errorf("unable to find user: %v", err)
	} else if err == user.ErrNotFound {
		return user.ErrNotFound
	}

	_, err = p.DBMap.Delete(u)
	if err != nil {
		return fmt.Errorf("unable to delete user: %v", err)
	}

	return nil
}

// Update only fiew fields from user
func (p *Provider) Update(u *user.User, fields map[string]interface{}) (err error) {
	var (
		sets     string
		args     []interface{}
		queryIdx = 0
	)
	for key, value := range fields {
		sets += " " + key + "=" + p.DBMap.Dialect.BindVar(queryIdx)
		args = append(args, value)
		queryIdx++
	}

	args = append(args, u.ID)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%s", // nolint: gas
		p.UserTableName, sets, p.DBMap.Dialect.BindVar(queryIdx))
	if _, err = p.DBMap.Exec(query, args...); err != nil {
		return fmt.Errorf("unable to update user informations: %v", err)
	}
	return nil
}
