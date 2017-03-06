package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/krostar/nebulo/config"
	"github.com/krostar/nebulo/router/httperror"
	"github.com/krostar/nebulo/tools/cert"
	"github.com/labstack/echo"
)

// UserCreate handle the route /user/.
// Return a CRT generated from the CRS submitted and the CA.
/**
 * @api {post} /user/ Register an account
 * @apiDescription Create a user and allow him to connect to restricted areas of the API
 * @apiName Register
 * @apiGroup User
 *
 * @apiExample {curl} Usage example
 *		$>curl "http://127.0.0.1:17241/version"
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
 * @apiError (Errors 4XX) {json} 401 Unauthorized
 * @apiError (Errors 4XX) {json} 404 Bad Request
 * @apiError (Errors 5XX) {json} 500 Internal server error
*/
func UserCreate(c echo.Context) error {
	// TODO: look for existing user with same public key

	// load certificate authority
	caCRT, caPrivateKey, err := cert.CAFromFiles(
		config.Config.TLSClientsCACertFile,
		config.Config.TLSClientsCAKeyFile,
		[]byte(config.Config.TLSClientsCAKeyPassword),
	)
	if err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to fetch certificate authority files: %v", err))
	}

	// get client certificate request from body
	bodyLength, err := strconv.ParseInt(c.Request().Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return httperror.HTTPBadRequestError(fmt.Errorf("bad content-length: %v", err))
	}
	rawBodyReader := bytes.NewBuffer(make([]byte, 0, bodyLength))
	_, err = rawBodyReader.ReadFrom(c.Request().Body)
	if err != nil {
		return httperror.HTTPBadRequestError(fmt.Errorf("unable to read from raw body: %v", err))
	}
	rawBody := rawBodyReader.Bytes()
	clientCSR, err := cert.CSRFromPEM(rawBody, nil)
	if err != nil {
		return httperror.HTTPBadRequestError(fmt.Errorf("unable to convert PEM certificate to certificate request: %v", err))
	}

	// create client certificate template
	clientCRTTemplate := x509.Certificate{
		Signature:          clientCSR.Signature,
		SignatureAlgorithm: clientCSR.SignatureAlgorithm,
		PublicKeyAlgorithm: clientCSR.PublicKeyAlgorithm,
		PublicKey:          clientCSR.PublicKey,
		SerialNumber:       big.NewInt(2),
		Issuer:             caCRT.Subject,
		Subject:            clientCSR.Subject,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(24 * time.Hour),
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// create client certificate from template and CA public key
	clientCRTRaw, err := x509.CreateCertificate(rand.Reader, &clientCRTTemplate, caCRT, clientCSR.PublicKey, caPrivateKey)
	if err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to create certificate: %v", err))
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
