package router

import (
	"crypto/tls"
	"net"
	"strconv"

	"github.com/krostar/nebulo/config"
	"github.com/krostar/nebulo/env"
	njwt "github.com/krostar/nebulo/router/auth/jwt"
	"github.com/krostar/nebulo/router/handler"
	"github.com/krostar/nebulo/router/httperror"
	nmiddleware "github.com/krostar/nebulo/router/middleware"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	if env.Environment(config.Config.Environment) == env.DEV {
		router.Debug = true
	} else {
		router.Debug = false
	}
	router.HTTPErrorHandler = httperror.ErrorHandler

	router.Server.Addr = environment.Address + ":" + strconv.Itoa(environment.Port)

	setupMiddlewares()
	setupRoutes()

}

func createTLSConfig(certFile string, keyFile string) (config *tls.Config, err error) {
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	config = &tls.Config{
		Certificates: []tls.Certificate{cer},
		//     RootCAs *x509.CertPool
		//     ServerName string
		// ClientAuth: tls.RequireAndVerifyClientCert,
		//     ClientCAs *x509.CertPool
		//     InsecureSkipVerify bool
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   false,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384},
	}

	// w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	return config, nil
}

func setupMiddlewares() {
	router.Use(nmiddleware.Log())
	router.Use(nmiddleware.Recover()) // in case of panic, recover and don't quit
	router.Use(middleware.RemoveTrailingSlash())
	router.Use(nmiddleware.Headers())

	JWTConfig := middleware.JWTConfig{
		Claims:        &njwt.Claims{},
		SigningKey:    []byte(config.Config.JWTSecret),
		AuthScheme:    "Bearer",
		SigningMethod: "HS512",
		ContextKey:    "jwt",
	}
	puMdw["auth"] = middleware.JWTWithConfig(JWTConfig)
}

func setupRoutes() {
	router.GET("/version", handler.Version)

	// domain/auth/...
	auth := router.Group("/auth")
	auth.POST("/login", handler.AuthLogin)
	auth.POST("/login/verify", handler.AuthLoginVerify, puMdw["auth"])
	auth.GET("/logout", handler.AuthLogout, puMdw["auth"])

	// domain/user/...
	user := router.Group("/user")
	user.GET("/:user", handler.UserInfos, puMdw["auth"]) //user profile infos
	// user.POST("/", handler.UserCreate)                 //add user profile
	// user.PUT("/:user", handler.UserEdit, mdwAuth)      //edit user profile
	// user.DELETE("/:user", handler.UserDelete, mdwAuth) //edit user profile
	//
	// // domain/chans
	// router.GET("/chans", handler.ChansList, mdwAuth) //list all channels
	//
	// // domain/chan/...
	// // identified users are required to make these calls
	// //     that's why everything using channel group use auth middleware
	// channel := router.Group("/chan", mdwAuth)
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
	// message.DELETE("/:message", handler.ChanMessageDelete) //delete a specific message
}

func run(environment *env.Config, tlsConfig *tls.Config) error {
	setupRouter(environment)
	router.Server.Addr = environment.Address + ":" + strconv.Itoa(environment.Port)

	listener, err := net.Listen("tcp4", router.Server.Addr)
	if err != nil {
		return err
	}

	if tlsConfig != nil {
		router.Server.TLSConfig = tlsConfig
		listener = tls.NewListener(listener, router.Server.TLSConfig)
	}
	return router.Server.Serve(listener)
}

// Run start the routeur
func Run(environment *env.Config) error {
	return run(environment, nil)
}

// RunTLS start the routeur and use encryption to communicate
func RunTLS(environment *env.Config, certFile string, keyFile string) error {
	tlsConfig, err := createTLSConfig(certFile, keyFile)
	if err != nil {
		return err
	}
	return run(environment, tlsConfig)
}
