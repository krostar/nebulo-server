package cert

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
)

// CAFromFiles load files from path and return a certificate and a password key
func CAFromFiles(certPath string, keyPath string, keyPassword []byte) (cert *x509.Certificate, key crypto.Signer, err error) {
	cert, err = ParseCertificatePEMFromFile(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to decode PEM encoded CA certificate file: %v", err)
	}

	caPrivateKey, err := ParsePrivateKeyPEMFromFile(keyPath, keyPassword)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to decode PEM encoded CA private key file: %v", err)
	}
	return cert, caPrivateKey, err
}

// Revoke the certificate made him unusable for login
func Revoke(cert *x509.Certificate) (err error) {
	// this is a temporary solution to avoid having a bigger infrastructure with an OCSP server
	f, err := os.OpenFile("certs/revokated.info", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("unable to open revokated certs file: %v", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err = f.WriteString(base64.RawStdEncoding.EncodeToString(cert.Raw) + "\n"); err != nil {
		return fmt.Errorf("unable to open revokated certs file: %v", err)
	}
	return nil
}
