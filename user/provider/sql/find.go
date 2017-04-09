package sql

import (
	"crypto/x509"
	"fmt"

	"github.com/krostar/nebulo/user"
)

// FindByPublicKey is used to find a user from his public key
func (p *Provider) FindByPublicKey(publicKeyAlgo x509.PublicKeyAlgorithm, publicKey interface{}) (u *user.User, err error) {
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal public key: %v", err)
	}

	return p.find("key_public_der", publicKeyDER)
}

// FindByID is used to find a user from his ID
func (p *Provider) FindByID(id int) (u *user.User, err error) {
	return p.find("id", id)
}

func (p *Provider) find(field string, value interface{}) (u *user.User, err error) {
	u = new(user.User)

	if p.DB.Where(field+" = ?", value).First(u).RecordNotFound() {
		return nil, user.ErrNotFound
	}
	if err = p.DB.Error; err != nil {
		return nil, fmt.Errorf("unable to select user in db: %v", err)
	}

	return u, nil
}
