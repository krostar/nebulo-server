package handler

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/krostar/nebulo-golib/router/httperror"
	"github.com/krostar/nebulo-golib/tools/cert"
	"github.com/krostar/nebulo-server/config"
	"github.com/krostar/nebulo-server/user"
	up "github.com/krostar/nebulo-server/user/provider"
	"github.com/labstack/echo"
)

// UserCreate handle the route POST /user/.
// Return a CRT generated from the CRS submitted and the CA.
/**
 * @api {post} /user/ Register an account
 * @apiDescription Create a user and allow him to connect to restricted areas of the API
 * @apiName User - Create profile
 * @apiGroup User
 *
 * @apiExample {curl} Usage example
 *		$>curl -X POST --cacert ca.crt -v "https://api.nebulo.io/user/" --data-binary "@user.csr"
 *
 * @apiSuccess (Success) {nothing} 201 Created
 * @apiSuccessExample {binary} Success example
 *		HTTP/1.1 201 "Created"
 *		-----BEGIN CERTIFICATE-----
		MIIE6jCCAtKgAwIBAgIBAjANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZuZWJ1
		bG8wHhcNMTcwMzA4MTIxMDIxWhcNMTcwMzA5MTIxMDIxWjAWMRQwEgYDVQQDEwtK
		...
		93FwQ9M4vipScDcrkyj9X9vueWzv7GBK2npXXsXoAVecLkLL5P6MMi8z7wcmlUSB
		FG4WG+sgP5x/bNY5fZ4=
		-----END CERTIFICATE-----
 *
 * @apiError (Errors 4XX) {json} 400 Bad Request: unable to load user certificate request
 * @apiError (Errors 4XX) {json} 409 Conflict: user already exist
 * @apiError (Errors 5XX) {json} 500 Internal server error: server failed to handle the request
*/
func UserCreate(c echo.Context) (err error) {
	// create client certificate template
	clientCSR, clientCRTRaw, err := signCertificate(c.Request().Body, c.Request().Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	// create a user with this public key
	storablePublicKey, err := x509.MarshalPKIXPublicKey(clientCSR.PublicKey)
	if err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to marshal public key: %v", err))
	}
	newUser := &user.User{
		PublicKeyDER: storablePublicKey,
		FingerPrint:  cert.FingerprintSHA256(storablePublicKey),
	}
	if _, err = up.P.Create(newUser); err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to register user in user provider: %v", err))
	}

	// send back the generated certificate
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Header().Add("Content-Type", "application/x-x509-user-cert")
	err = pem.Encode(c.Response().Writer, &pem.Block{Type: "CERTIFICATE", Bytes: clientCRTRaw})
	if err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to send back the certificate: %v", err))
	}

	return nil
}

func signCertificate(requestBody io.Reader, contentLengthHeader string) (clientCSR *x509.CertificateRequest, clientCRTRaw []byte, err error) {
	// load required certificate
	clientCSR, caCert, caPrivateKey, err := loadCertificate(requestBody, contentLengthHeader)
	if err != nil {
		return nil, nil, err
	}

	// check if a user exist with this public key
	if _, err = up.P.FindByPublicKey(clientCSR.PublicKey); err == nil {
		return nil, nil, httperror.UserExist()
	} else if err != nil && err != user.ErrNotFound {
		return nil, nil, httperror.HTTPInternalServerError(err)
	}

	clientCRTTemplate := x509.Certificate{
		Signature:          clientCSR.Signature,
		SignatureAlgorithm: clientCSR.SignatureAlgorithm,
		PublicKeyAlgorithm: clientCSR.PublicKeyAlgorithm,
		PublicKey:          clientCSR.PublicKey,
		SerialNumber:       big.NewInt(2),
		Issuer:             caCert.Subject,
		Subject:            clientCSR.Subject,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(7 * time.Hour * 24),
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// create/sign the request with the client CA
	clientCRTRaw, err = x509.CreateCertificate(rand.Reader, &clientCRTTemplate, caCert, clientCSR.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, httperror.HTTPInternalServerError(fmt.Errorf("unable to create certificate: %v", err))
	}
	return clientCSR, clientCRTRaw, nil
}

// get client certificate request from body
func loadCertificate(requestBody io.Reader, contentLengthHeader string) (clientCSR *x509.CertificateRequest, caCert *x509.Certificate, caPrivateKey crypto.PrivateKey, err error) {
	// check body
	bodyLength, err := strconv.ParseInt(contentLengthHeader, 10, 64)
	if err != nil {
		return nil, nil, nil, httperror.HTTPBadRequestError(fmt.Errorf("bad content-length: %v", err))
	}
	if bodyLength < 210 {
		return nil, nil, nil, httperror.HTTPBadRequestError(errors.New("no csr submitted"))
	}

	// parse body to get certificate request
	rawBodyReader := bytes.NewBuffer(make([]byte, 0, bodyLength))
	_, err = rawBodyReader.ReadFrom(requestBody)
	if err != nil {
		return nil, nil, nil, httperror.HTTPBadRequestError(fmt.Errorf("unable to read from raw body: %v", err))
	}
	clientCSR, err = cert.ParseCSR(rawBodyReader.Bytes())
	if err != nil {
		return nil, nil, nil, httperror.HTTPBadRequestError(fmt.Errorf("unable to convert raw body to certificate request: %v", err))
	}

	// load certificate authority
	caCert, caPrivateKey, err = cert.CertAndKeyFromFiles(
		config.Config.Run.TLS.ClientsCA.Cert,
		config.Config.Run.TLS.ClientsCA.Key,
		[]byte(config.Config.Run.TLS.ClientsCA.KeyPassword),
	)
	if err != nil {
		return nil, nil, nil, httperror.HTTPInternalServerError(fmt.Errorf("unable to fetch certificate authority files: %v", err))
	}
	return clientCSR, caCert, caPrivateKey, nil
}
