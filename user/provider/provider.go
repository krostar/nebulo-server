package provider

import (
	gp "github.com/krostar/nebulo-golib/provider"
	"github.com/krostar/nebulo-server/user"
)

// Provider contains all the methods needed to manage users
type Provider interface {
	gp.TablesManagement

	Login(u *user.User) (err error)
	Create(userToAdd *user.User) (u *user.User, err error)
	Delete(u *user.User) (err error)

	FindByPublicKey(publicKey interface{}) (u *user.User, err error)
	FindByPublicKeyDER(publicKeyDER []byte) (u *user.User, err error)
	FindByPublicKeyDERBase64(publicKeyDERBase64 string) (u *user.User, err error)
	FindByID(ID int) (u *user.User, err error)

	Update(u *user.User, fields map[string]interface{}) (err error)
}

// P is the selected provider
var P Provider
