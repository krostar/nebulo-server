package user

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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
	ID           int       `json:"-" gorm:"column:id; primary_key; not null"`
	PublicKeyDER []byte    `json:"-" gorm:"column:key_public_der; size:2000; not null"`
	FingerPrint  string    `json:"key_fingerprint" gorm:"column:key_fingerprint; size:51; not null"`
	DisplayName  string    `json:"display_name" gorm:"column:display_name; size:42" validator-update:"string=length:max|42"`
	Signup       time.Time `json:"signup" gorm:"column:signup; not null" sql:"DEFAULT:current_timestamp"`
	LoginFirst   time.Time `json:"login_first" gorm:"column:login_first; not null" sql:"DEFAULT:'1970-01-01 00:00:00'"`
	LoginLast    time.Time `json:"login_last" gorm:"column:login_last; not null" sql:"DEFAULT:'1970-01-01 00:00:00'"`
}

// MarshalJSON overload the default user json marshal to add fields useful for clients
func (u *User) MarshalJSON() ([]byte, error) {
	// create a new type to avoid infinite Marshal recursion
	type fakeUser User
	type userJSON struct {
		*fakeUser
		PublicKeyDERBase64 string `json:"public_key_der_b64"`
	}

	mUser := &userJSON{
		fakeUser:           (*fakeUser)(u),
		PublicKeyDERBase64: base64.StdEncoding.EncodeToString(u.PublicKeyDER),
	}

	return json.Marshal(mUser)
}

// Repr return an uniq representation of a given user this function should nerver
// be used to find a user, only to differentiate it visually because of possible collision
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
