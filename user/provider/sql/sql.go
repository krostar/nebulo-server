package sql

import (
	"errors"
	"fmt"
	"time"

	gp "github.com/krostar/nebulo-golib/provider"
	"github.com/krostar/nebulo-server/channel/provider"
	"github.com/krostar/nebulo-server/user"
	"github.com/labstack/gommon/log"
)

// Provider implements the methods needed to manage users
// for every SQL based database
type Provider struct {
	*gp.RootProvider
	provider.Provider
}

// Login update field on user login
func (p *Provider) Login(u *user.User) (err error) {
	if u == nil {
		return user.ErrNil
	}

	now := time.Now()

	updates := user.User{LoginLast: now.UTC()}
	if u.LoginFirst.UTC() == time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC) {
		log.Debugf("updates.loginFirst is zero: %v", u.LoginFirst)
		updates.LoginFirst = now
	}
	if err = p.DB.Model(u).Updates(updates).Error; err != nil {
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
	_, err = p.FindByPublicKeyDER(userToAdd.PublicKeyDER)
	if err != nil && err != user.ErrNotFound {
		return nil, fmt.Errorf("unable to find user: %v", err)
	}
	if err != user.ErrNotFound {
		return nil, errors.New("an user already exist with this public key")
	}

	if err = p.DB.Create(userToAdd).Error; err != nil {
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

	if err = p.DB.Delete(u).Error; err != nil {
		return fmt.Errorf("unable to delete user: %v", err)
	}

	return nil
}

// Update only fiew fields from user
func (p *Provider) Update(u *user.User, fields map[string]interface{}) (err error) {
	if u == nil {
		return user.ErrNil
	}

	if err = p.DB.Model(u).Updates(fields).Error; err != nil {
		return fmt.Errorf("unable to update user informations: %v", err)
	}
	return nil
}
