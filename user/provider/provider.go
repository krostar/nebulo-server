package provider

import (
	"crypto/x509"

	"github.com/krostar/nebulo/user"
)

// Provider represent the way we can interact with a provider
// to get informations about a user
type Provider interface {
	Register(u *user.User) (err error)
	FindByPublicKey(publicKeyAlgo x509.PublicKeyAlgorithm, publicKey interface{}) (u *user.User, err error)
	FindByID(ID int) (u *user.User, err error)
	Save(u *user.User) (err error)
}

// P is the currently used provider
var P Provider

// Use set the new provider as the provider to use
func Use(newProvider Provider) (err error) {
	P = newProvider
	return nil
}
