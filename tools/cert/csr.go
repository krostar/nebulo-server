package cert

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
)

// ParseCSR parses a PEM-encoded PKCS #10 certificate signing request.
func ParseCSR(in []byte) (csr *x509.CertificateRequest, err error) {
	in = bytes.TrimSpace(in)
	block, _ := pem.Decode(in)
	if block == nil {
		return nil, errors.New("empty pem block")
	}
	if block.Type != "NEW CERTIFICATE REQUEST" && block.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New("bad certificate type")
	}

	csr, err = x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate request: %v", err)
	}

	err = CheckCSRSignature(csr, csr.SignatureAlgorithm, csr.RawTBSCertificateRequest, csr.Signature)
	if err != nil {
		return nil, fmt.Errorf("bad signature for csr: %v", err)
	}

	return csr, nil
}

// CheckCSRSignature verifies a signature made by the key on a CSR.
// nolint: gocyclo
func CheckCSRSignature(csr *x509.CertificateRequest, algo x509.SignatureAlgorithm, signed, signature []byte) error {
	var hashType crypto.Hash
	switch algo {
	case x509.SHA1WithRSA, x509.ECDSAWithSHA1:
		hashType = crypto.SHA1
	case x509.SHA256WithRSA, x509.ECDSAWithSHA256:
		hashType = crypto.SHA256
	case x509.SHA384WithRSA, x509.ECDSAWithSHA384:
		hashType = crypto.SHA384
	case x509.SHA512WithRSA, x509.ECDSAWithSHA512:
		hashType = crypto.SHA512
	default:
		return x509.ErrUnsupportedAlgorithm
	}
	if !hashType.Available() {
		return x509.ErrUnsupportedAlgorithm
	}

	h := hashType.New()
	if _, err := h.Write(signed); err != nil {
		return err
	}
	digest := h.Sum(nil)

	switch pub := csr.PublicKey.(type) {
	case *rsa.PublicKey:
		return rsa.VerifyPKCS1v15(pub, hashType, digest, signature)
	case *ecdsa.PublicKey:
		ecdsaSig := new(struct{ R, S *big.Int })
		_, err := asn1.Unmarshal(signature, ecdsaSig)
		if err == nil && ecdsaSig.R.Sign() <= 0 || ecdsaSig.S.Sign() <= 0 {
			err = errors.New("x509: ECDSA signature contained zero or negative values")
		}
		if err == nil && !ecdsa.Verify(pub, digest, ecdsaSig.R, ecdsaSig.S) {
			err = errors.New("x509: ECDSA verification failure")
		}
		return err
	default:
		return errors.New("unknow public key type")
	}
}
