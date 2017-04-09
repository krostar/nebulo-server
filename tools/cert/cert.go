package cert

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// EncodeCertificatePEM encodes a single x509 certficates to PEM
func EncodeCertificatePEM(cert *x509.Certificate) (encoded []byte, err error) {
	var buffer bytes.Buffer
	if err = pem.Encode(&buffer, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}); err != nil {
		return nil, fmt.Errorf("unable to encode pem certificate: %v", err)
	}
	return buffer.Bytes(), nil
}

// ParseCertificatePEMFromFile call ParseCertificatePEM with the content of a file
func ParseCertificatePEMFromFile(filename string) (cert *x509.Certificate, err error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %v", filename, err)
	}
	return ParseCertificatePEM(raw)
}

// ParseCertificatePEM parses and returns a PEM-encoded certificate
func ParseCertificatePEM(certPEM []byte) (cert *x509.Certificate, err error) {
	certPEM = bytes.TrimSpace(certPEM)
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, errors.New("empty pem block")
	}
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate: %v", err)
	}
	return cert, nil
}

// ParsePrivateKeyPEMFromFile call ParsePrivateKeyPEM with the content of a file
func ParsePrivateKeyPEMFromFile(filename string, password []byte) (key crypto.Signer, err error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %v", filename, err)
	}
	return ParsePrivateKeyPEM(raw, password)
}

// ParsePrivateKeyPEM parses and returns a PEM-encoded private key.
// The private key may be a potentially encrypted PKCS#8, PKCS#1, or elliptic private key.
// nolint: gocyclo
func ParsePrivateKeyPEM(keyPEM []byte, password []byte) (key crypto.Signer, err error) {
	keyPEM = bytes.TrimSpace(keyPEM)
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, errors.New("empty pem block")
	}
	if procType, ok := block.Headers["Proc-Type"]; ok && strings.Contains(procType, "ENCRYPTED") && password != nil {
		if block.Bytes, err = x509.DecryptPEMBlock(block, password); err != nil {
			return nil, fmt.Errorf("unable to decrypt pem block: %v", err)
		}
	}

	var privateKey interface{}
	privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			privateKey, err = x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				return nil, errors.New("unable to parse private key: every parsing function failed")
			}
		}
	}

	switch privateKey.(type) {
	case *rsa.PrivateKey:
		return privateKey.(*rsa.PrivateKey), nil
	case *ecdsa.PrivateKey:
		return privateKey.(*ecdsa.PrivateKey), nil
	default:
		return nil, errors.New("unknow private key type")
	}
}

// VerifyCertificate check for a certificate revokation
func VerifyCertificate(cert *x509.Certificate) (revoked bool, err error) {
	if !time.Now().Before(cert.NotAfter) {
		return true, nil
	} else if !time.Now().After(cert.NotBefore) {
		return true, nil
	}

	// this is a temporary solution to avoid having a bigger infrastructure with an OCSP server
	f, err := os.OpenFile("certs/revokated.info", os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return false, fmt.Errorf("unable to open revokated certs file: %v", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	check := base64.RawStdEncoding.EncodeToString(cert.Raw)

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == check {
			return true, nil
		}
	}

	return false, nil
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

// PublicKeyAlgorithmToString return a string representation
// of a x509.PublicKeyAlgorithm
func PublicKeyAlgorithmToString(pka x509.PublicKeyAlgorithm) string {
	switch pka {
	case x509.RSA:
		return "RSA"
	case x509.DSA:
		return "DSA"
	case x509.ECDSA:
		return "ECDSA"
	default:
		return "unknown"
	}
}
