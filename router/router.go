package router

import (
	"strconv"
	"time"

	graceful "gopkg.in/tylerb/graceful.v1"

	"github.com/krostar/nebulo/config"
	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
	njwt "github.com/krostar/nebulo/router/auth/jwt"
	"github.com/krostar/nebulo/router/handler"
	"github.com/krostar/nebulo/router/httperror"
	nmiddleware "github.com/krostar/nebulo/router/middleware"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	router *echo.Echo
)

// init define some useful-always-used parameters to echo.Echo router
func init() {
	router = echo.New()
}

func setupRouter() {
	if env.Environment(config.Config.Environment) == env.DEV {
		router.Debug = true
	} else {
		router.Debug = false
	}

	router.HTTPErrorHandler = httperror.ErrorHandler
	router.Use(nmiddleware.Log())
	router.Use(nmiddleware.Recover()) // in case of panic, recover and don't quit
	router.Use(middleware.RemoveTrailingSlash())

	JWTConfig := middleware.JWTConfig{
		Claims:        &njwt.Claims{},
		SigningKey:    []byte(config.Config.JWTSecret),
		AuthScheme:    "Bearer",
		SigningMethod: "HS512",
		ContextKey:    "jwt",
	}
	mdwAuth := middleware.JWTWithConfig(JWTConfig)

	router.GET("/version", handler.Version)

	// domain/auth/...
	auth := router.Group("/auth")
	auth.POST("/login", handler.AuthLogin)
	auth.POST("/login/verify", handler.AuthLoginVerify, mdwAuth)
	auth.GET("/logout", handler.AuthLogout, mdwAuth)

	// domain/user/...
	user := router.Group("/user")
	user.GET("/:user", handler.UserInfos, mdwAuth) //user profile infos
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

// Run start the routeur
func Run(environment *env.Config) error {
	setupRouter()
	router.Server.Addr = environment.Address + ":" + strconv.Itoa(environment.Port)
	log.Infoln("Starting router on", router.Server.Addr)
	return graceful.ListenAndServe(router.Server, 10*time.Second)
}

// RunTLS start the routeur and use encryption to communicate
func RunTLS(environment *env.Config, certFile string, keyFile string) error {
	setupRouter()
	router.Server.Addr = environment.Address + ":" + strconv.Itoa(environment.Port)
	log.Infoln("Starting router on", router.Server.Addr)
	return graceful.ListenAndServeTLS(router.Server, certFile, keyFile, 10*time.Second)
}
