package cert

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

func CAFromFiles(crtPath string, keyPath string, keyPassword []byte) (crt *x509.Certificate, key *rsa.PrivateKey, err error) {
	caCRTRaw, err := DERFromPEMFile(crtPath, nil)
	if err != nil {
		return nil, nil, err
	}
	caCRT, err := x509.ParseCertificate(caCRTRaw)
	if err != nil {
		return nil, nil, err
	}

	caPrivateKeyRaw, err := DERFromPEMFile(keyPath, keyPassword)
	if err != nil {
		return nil, nil, err
	}
	caPrivateKey, err := x509.ParsePKCS1PrivateKey(caPrivateKeyRaw)
	if err != nil {
		return nil, nil, err
	}
	return caCRT, caPrivateKey, err
}

func DERFromPEMFile(filename string, password []byte) (der []byte, err error) {
	rawCert, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return LoadPEM(rawCert, password)
}

func CSRFromPEM(rawCSR []byte, password []byte) (clientCSR *x509.CertificateRequest, err error) {
	clientCSRRAW, err := LoadPEM(rawCSR, nil)
	if err != nil {
		return nil, err
	}
	clientCSR, err = x509.ParseCertificateRequest(clientCSRRAW)
	if err != nil {
		return nil, err
	}
	if err = clientCSR.CheckSignature(); err != nil {
		return nil, err
	}
	return clientCSR, nil
}

func LoadPEM(rawCertificate []byte, password []byte) (der []byte, err error) {
	pemBlock, _ := pem.Decode(rawCertificate)
	if pemBlock == nil {
		return nil, errors.New("pem.Decode failed")
	}

	der = pemBlock.Bytes
	if password != nil {
		der, err = x509.DecryptPEMBlock(pemBlock, password)
		if err != nil {
			return nil, err
		}
	}
	return der, nil
}
