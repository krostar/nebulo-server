package provider

import (
	"crypto/x509"

	gp "github.com/krostar/nebulo/provider"
	"github.com/krostar/nebulo/user"
)

// Provider contains all the methods needed to manage users
type Provider interface {
	gp.TablesManagement

	Login(u *user.User) (err error)
	Create(userToAdd *user.User) (u *user.User, err error)
	Delete(u *user.User) (err error)

	FindByPublicKey(publicKeyAlgo x509.PublicKeyAlgorithm, publicKey interface{}) (u *user.User, err error)
	FindByID(ID int) (u *user.User, err error)

	Update(u *user.User, fields map[string]interface{}) (err error)
}

// P is the selected provider
var P Provider
