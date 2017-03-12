package file

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/krostar/nebulo/user"
	"github.com/krostar/nebulo/user/provider"

	validator "gopkg.in/validator.v2"
)

// ErrProviderNil is throw when the file provider object is nil
var ErrProviderNil = errors.New("user file provider is nil")

// Config is the configuration for the user provider file
type Config struct {
	Filepath string `validate:"file=readable:createifmissing"`
}

// ProviderFile implement Provider and contain the configuration
// and the user list loaded from the file
type ProviderFile struct {
	provider.Provider
	config *Config
	users  []*user.User
}

// NewFromConfig return a new ProviderFile based on the configuration
func NewFromConfig(config interface{}) (p *ProviderFile, err error) {
	newConfig, ok := config.(*Config)
	if !ok {
		return nil, errors.New("unable to cast config to *file.Config")
	}
	if err = validator.Validate(config); err != nil {
		return nil, fmt.Errorf("user file provider configuration validation failed: %v", err)
	}

	var users []*user.User

	// load users from file
	raw, err := ioutil.ReadFile(newConfig.Filepath)
	if err != nil {
		return nil, fmt.Errorf("user provider file loading failed: %v", err)
	}

	if len(raw) != 0 {
		if err = json.Unmarshal(raw, &users); err != nil {
			return nil, fmt.Errorf("user provider file parsing failed: %v", err)
		}
	}

	return &ProviderFile{
		config: newConfig,
		users:  users,
	}, nil
}

// Register create a new User in the file
func (pf *ProviderFile) Register(u *user.User) error {
	if pf == nil {
		return ErrProviderNil
	}

	maxID := 0
	for _, u := range pf.users {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	u.ID = maxID + 1

	pf.users = append(pf.users, u)
	return pf.Save(u)
}

// FindByID is used to find a user from his ID
func (pf *ProviderFile) FindByID(ID int) (u *user.User, err error) {
	if pf == nil {
		return nil, ErrProviderNil
	}

	for _, u := range pf.users {
		if u.ID == ID {
			return u, nil
		}
	}
	return nil, user.ErrNotFound
}

// FindByPublicKey is used to find a user from his public key
func (pf *ProviderFile) FindByPublicKey(publicKeyAlgo x509.PublicKeyAlgorithm, publicKey interface{}) (u *user.User, err error) {
	if pf == nil {
		return nil, ErrProviderNil
	}

	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal public key: %v", err)
	}
	for _, u := range pf.users {
		if u.PublicKeyAlgorithm == publicKeyAlgo && bytes.Equal(u.PublicKeyDER, publicKeyDER) {
			return u, nil
		}
	}
	return nil, user.ErrNotFound
}

// Save save modifications made on user struct
func (pf *ProviderFile) Save(u *user.User) (err error) {
	if pf == nil {
		return ErrProviderNil
	}

	raw, err := json.MarshalIndent(pf.users, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pf.config.Filepath, raw, 640)
}
