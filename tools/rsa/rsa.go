package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

var (
	// ErrPemDecodeFailed is an error returned
	// when something went wrong with pem key decoding
	ErrPemDecodeFailed = errors.New("unable to decode pem key")

	// ErrRSAPublicKeyCastFailed is an error returned
	// when something went wrong with rsa public key casting
	ErrRSAPublicKeyCastFailed = errors.New("failed to cast public key to RSA public key")
)

// LoadPublicKeyBase64 convert a base64 representation of an rsa public key to the *rsa.PublicKey type
func LoadPublicKeyBase64(b64PublicKey string) (*rsa.PublicKey, error) {
	publicKeyRepr, err := base64.StdEncoding.DecodeString(b64PublicKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyRepr)
	if block == nil {
		return nil, ErrPemDecodeFailed
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrRSAPublicKeyCastFailed
	}

	return rsaPublicKey, nil
}
