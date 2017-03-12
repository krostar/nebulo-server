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
	// ErrUserNil is throw when an user is nil
	ErrUserNil = errors.New("user is nil")
)

// User is the modelisation of an user
type User struct {
	ID                 int                     `json:"id"`
	PublicKeyDER       []byte                  `json:"key_public_der"`
	PublicKeyAlgorithm x509.PublicKeyAlgorithm `json:"key_public_algo"`
	FingerPrint        string                  `json:"key_fingerprint"`
	DisplayName        string                  `json:"display_name"`
	SignUp             time.Time               `json:"signup"`
	LoginFirst         time.Time               `json:"login_first"`
	LoginLast          time.Time               `json:"login_last"`
}

// Repr return an uniq representation of a given user
func (u *User) Repr() (string, error) {
	if u == nil {
		return "", ErrUserNil
	}

	userRepr := fmt.Sprintf("%s - %d", u.FingerPrint, u.ID)

	hash := sha256.New()
	if _, err := hash.Write([]byte(userRepr)); err != nil {
		return "", fmt.Errorf("unable to create sha256 hash of user representation: %v", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
