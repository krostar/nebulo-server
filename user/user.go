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
	ID                 int                     `json:"-" gorm:"column:id; primary_key; not null"`
	PublicKeyDER       []byte                  `json:"-" gorm:"column:key_public_der; size:2000; not null"`
	PublicKeyAlgorithm x509.PublicKeyAlgorithm `json:"-" gorm:"column:key_public_algorithm; type:tinyint(1); not null"`
	FingerPrint        string                  `json:"key_fingerprint" gorm:"column:key_fingerprint; size:51; not null"`
	DisplayName        string                  `json:"display_name" gorm:"column:display_name; size:42" validator-update:"string=length:max|42"`
	Signup             time.Time               `json:"signup" gorm:"column:signup; not null" sql:"DEFAULT:current_timestamp"`
	LoginFirst         time.Time               `json:"login_first" gorm:"column:login_first; not null" sql:"DEFAULT:'1970-01-01 00:00:00'"`
	LoginLast          time.Time               `json:"login_last" gorm:"column:login_last; not null" sql:"DEFAULT:'1970-01-01 00:00:00'"`
}

// Repr return an uniq representation of a given user
func (u *User) Repr() (string, error) {
	if u == nil {
		return "", ErrNil
	}

	// FingerPrint is not enough because of possible collision, anyway Repr should never
	// be used to find a user, only to differentiate visually

	userRepr := fmt.Sprintf("%s - %d", u.FingerPrint, u.ID)

	hash := sha256.New()
	if _, err := hash.Write([]byte(userRepr)); err != nil {
		return "", fmt.Errorf("unable to create sha256 hash of user representation: %v", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
