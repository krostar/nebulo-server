package router

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"

	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/router/handler"
	"github.com/krostar/nebulo/router/httperror"
	nmiddleware "github.com/krostar/nebulo/router/middleware"
	"github.com/labstack/echo"
)

var (
	router *echo.Echo
	puMdw  map[string]echo.MiddlewareFunc
)

// init define some useful-always-used parameters to echo.Echo router
func init() {
	router = echo.New()
	puMdw = make(map[string]echo.MiddlewareFunc)
}

func setupRouter(environment *env.Config) {
	if environment.Type == env.DEV {
		router.Debug = true
	} else {
		router.Debug = false
	}

	router.HTTPErrorHandler = httperror.ErrorHandler
	router.Server.Addr = environment.Address + ":" + strconv.Itoa(environment.Port)

	setupMiddlewares()
	setupRoutes()
}

func createTLSConfig(certFile string, keyFile string, clientCAFilepath string) (config *tls.Config, err error) {
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load tls key pair: %v", err)
	}

	// Load our CA certificate
	clientCAFile, err := ioutil.ReadFile(clientCAFilepath)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %v", clientCAFile, err)
	}

	clientCAPool := x509.NewCertPool()
	clientCAPool.AppendCertsFromPEM(clientCAFile)

	config = &tls.Config{
		Certificates: []tls.Certificate{cer},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ClientCAs:    clientCAPool,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   false,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384},
	}

	return config, nil
}

func setupMiddlewares() {
	router.Pre(nmiddleware.Recover()) // in case of panic, recover and don't quit
	router.Use(nmiddleware.Misc())
	router.Use(nmiddleware.Log())

	puMdw["auth"] = nmiddleware.Auth()
}

func setupRoutes() {
	router.GET("/version", handler.Version)

	// domain/user/...
	user := router.Group("/user")
	user.GET("/", handler.UserInfos, puMdw["auth"])     //user profile infos
	user.POST("/", handler.UserCreate)                  //create user profile
	user.PUT("/", handler.UserEdit, puMdw["auth"])      //edit user profile
	user.DELETE("/", handler.UserDelete, puMdw["auth"]) //delete user profile
	//
	// // domain/chans
	// router.GET("/chans", handler.ChansList, puMdw["auth"]) //list all channels
	//
	// // domain/chan/...
	// // identified users are required to make these calls
	// //     that's why everything using channel group use auth middleware
	// channel := router.Group("/chan", puMdw["auth"])
	// channel.GET("/:chan", handler.ChanInfos)     //get info for a specific channel
	// channel.POST("/", handler.ChanCreate)        //add a new channel
	// channel.PUT("/:chan", handler.ChanEdit)      //edit info of a specific channel
	// channel.DELETE("/:chan", handler.ChanDelete) //delete a specific channel
	//
	// // domain/chan/:chan/messages/...
	// messages := channel.Group("/:chan/messages")
	// messages.GET("/", handler.ChanMessagesList)      //get message list for a specific channel
	// messages.DELETE("/", handler.ChanMessagesDelete) //delete range of messages
	//
	// // domain/chan/:chan/message/...
	// message := channel.Group("/:chan/message")
	// message.POST("/", handler.ChanMessageCreate)           //add message to a specific channel
	// message.PUT("/:message", handler.ChanMessageEdit)      //edit a specific message
	// message.DELETE("/:message", handler.ChanMessa && config.Config.TLSKeyFile != ""geDelete) //delete a specific message
}

func run(environment *env.Config, tlsConfig *tls.Config) error {
	setupRouter(environment)
	router.Server.ErrorLog = log.New(ioutil.Discard, "", 0)
	router.Server.Addr = environment.Address + ":" + strconv.Itoa(environment.Port)

	listener, err := net.Listen("tcp4", router.Server.Addr)
	if err != nil {
		return fmt.Errorf("unable to listen on %s: %v", router.Server.Addr, err)
	}

	if tlsConfig != nil {
		router.Server.TLSConfig = tlsConfig
		listener = tls.NewListener(listener, router.Server.TLSConfig)
	}
	return router.Server.Serve(listener)
}

// RunTLS start the routeur and use encryption to communicate
func RunTLS(environment *env.Config, certFile string, keyFile string, clientsCAFile string) error {
	tlsConfig, err := createTLSConfig(certFile, keyFile, clientsCAFile)
	if err != nil {
		return fmt.Errorf("tls configuration error: %v", err)
	}
	return run(environment, tlsConfig)
}
