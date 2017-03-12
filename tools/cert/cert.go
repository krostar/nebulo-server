package cert

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

// CAFromFiles return a certificate object from a file
func CAFromFiles(crtPath string, keyPath string, keyPassword []byte) (crt *x509.Certificate, key *rsa.PrivateKey, err error) {
	caCRTRaw, err := DERFromPEMFile(crtPath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to decode PEM encoded CA certificate file: %v", err)
	}
	caCRT, err := x509.ParseCertificate(caCRTRaw)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse CA certificate: %v", err)
	}

	caPrivateKeyRaw, err := DERFromPEMFile(keyPath, keyPassword)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to decode PEM encoded CA private key file: %v", err)
	}
	caPrivateKey, err := x509.ParsePKCS1PrivateKey(caPrivateKeyRaw)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse CA key file: %v", err)
	}
	return caCRT, caPrivateKey, err
}

// DERFromPEMFile return a pem decoded byte array from a file
func DERFromPEMFile(filename string, password []byte) (der []byte, err error) {
	rawCert, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %v", filename, err)
	}
	return LoadPEM(rawCert, password)
}

// CSRFromPEM return a certificate request object from a pem decoded byte array
func CSRFromPEM(rawCSR []byte, password []byte) (clientCSR *x509.CertificateRequest, err error) {
	clientCSRRAW, err := LoadPEM(rawCSR, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to load PEM from raw certificate request: %v", err)
	}
	clientCSR, err = x509.ParseCertificateRequest(clientCSRRAW)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate request: %v", err)
	}
	if err = clientCSR.CheckSignature(); err != nil {
		return nil, fmt.Errorf("certificate request signature verification failed: %v", err)
	}
	return clientCSR, nil
}

// LoadPEM return a pem decoded byte array from a pem encoded byte array
func LoadPEM(rawCertificate []byte, password []byte) (der []byte, err error) {
	pemBlock, _ := pem.Decode(rawCertificate)
	if pemBlock == nil {
		return nil, errors.New("empty pem block")
	}

	der = pemBlock.Bytes
	if password != nil {
		der, err = x509.DecryptPEMBlock(pemBlock, password)
		if err != nil {
			return nil, fmt.Errorf("unable to decrypt pem block: %v", err)
		}
	}
	return der, nil
}

// FingerprintSHA256 returns the user presentation of the key's
// fingerprint as unpadded base64 encoded sha256 hash.
// This format was introduced from OpenSSH 6.8.
// https://www.openssh.com/txt/release-6.8
// https://tools.ietf.org/html/rfc4648#section-3.2 (unpadded base64 encoding)
// inspired from x/crypto/ssh package
func FingerprintSHA256(pubKeyDER []byte) string {
	sha256sum := sha256.Sum256(pubKeyDER)
	hash := base64.RawStdEncoding.EncodeToString(sha256sum[:])
	return "SHA256:" + hash
}
