package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
)

var (
	errUserNil = errors.New("user is nil")
)

// User is the modelisation of an user
type User struct {
	ID           int        `json:"id"`
	B64PublicKey string     `json:"key_rsa_public"`
	FingerPrint  string     `json:"key_fingerprint"`
	DisplayName  string     `json:"display_name"`
	SignUp       time.Time  `json:"signup"`
	LoginFirst   *time.Time `json:"login_first"`
	LoginLast    *time.Time `json:"login_last"`
}

// Repr return an uniq representation of a given user
func (u *User) Repr() (string, error) {
	if u == nil {
		return "", errUserNil
	}

	userRepr := fmt.Sprintf("%s - %d", u.FingerPrint, u.ID)

	hash := sha256.New()
	if _, err := hash.Write([]byte(userRepr)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
