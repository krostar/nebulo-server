package provider

import (
	"crypto/x509"
	"errors"

	"github.com/krostar/nebulo/user"
)

var (
	// SQLLogger is used to log every sql action
	SQLLogger *ORPLogger
)

func init() {
	SQLLogger = new(ORPLogger)
}

// Provider represent the way we can interact with a provider
// to get informations about a user
type Provider interface {
	Login(u *user.User) (err error)
	Register(userToAdd *user.User) (u *user.User, err error)
	FindByPublicKey(publicKeyAlgo x509.PublicKeyAlgorithm, publicKey interface{}) (u *user.User, err error)
	FindByID(ID int) (u *user.User, err error)
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
