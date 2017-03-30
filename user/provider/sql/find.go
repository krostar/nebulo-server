package sql

import (
	"crypto/x509"
	"database/sql"
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

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s=%s", // nolint: gas
		p.UserTableName, field, p.DBMap.Dialect.BindVar(0))
	if err = p.DBMap.SelectOne(u, query, value); err != nil {
		return nil, fmt.Errorf("unable to select user in db: %v", err)
	}
	if err == sql.ErrNoRows || u == nil {
		return nil, user.ErrNotFound
	}

	return u, nil
}
