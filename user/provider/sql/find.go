package sql

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/krostar/nebulo-server/user"
)

// FindByPublicKey is used to find a user from his public key
func (p *Provider) FindByPublicKey(publicKey interface{}) (u *user.User, err error) {
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal public key: %v", err)
	}
	return p.FindByPublicKeyDER(publicKeyDER)
}

// FindByPublicKeyDERBase64 is used to find a user from his public key der base64 formatted
func (p *Provider) FindByPublicKeyDERBase64(publicKeyDERBase64 string) (u *user.User, err error) {
	publicKeyDER, err := base64.StdEncoding.DecodeString(publicKeyDERBase64)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64 public key der: %v", err)
	}
	return p.FindByPublicKeyDER(publicKeyDER)
}

// FindByPublicKeyDER is used to find a user from his public key der formatted
func (p *Provider) FindByPublicKeyDER(publicKeyDER []byte) (u *user.User, err error) {
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
