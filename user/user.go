package user

import (
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNotFound is throw when an user is not found
	ErrNotFound = errors.New("user not found")
	// ErrNil is throw when an user is nil
	ErrNil = errors.New("user is nil")
)

// User is the modelisation of an user
type User struct {
	ID                 int                     `json:"-" db:"id, primarykey, autoincrement, notnull"`
	PublicKeyDER       []byte                  `json:"-" db:"key_public_der, notnull"`
	PublicKeyAlgorithm x509.PublicKeyAlgorithm `json:"-" db:"key_public_algo, notnull"`
	FingerPrint        string                  `json:"key_fingerprint" db:"key_fingerprint, size:51 notnull"`
	DisplayName        string                  `json:"display_name" db:"display_name, size:42" validator-update:"string=length:max|42"`
	SignUp             time.Time               `json:"signup" db:"signup, notnull"`
	LoginFirst         time.Time               `json:"login_first" db:"login_first"`
	LoginLast          time.Time               `json:"login_last" db:"login_last"`
}

// Repr return an uniq representation of a given user
func (u *User) Repr() (string, error) {
	if u == nil {
		return "", ErrNil
	}

	userRepr := fmt.Sprintf("%s - %d", u.FingerPrint, u.ID)

	hash := sha256.New()
	if _, err := hash.Write([]byte(userRepr)); err != nil {
		return "", fmt.Errorf("unable to create sha256 hash of user representation: %v", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
