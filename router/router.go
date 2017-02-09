package router

import (
	"strconv"
	"time"

	graceful "gopkg.in/tylerb/graceful.v1"

	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/handler"
	"github.com/krostar/nebulo/log"
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
	router.HTTPErrorHandler = httperror.ErrorHandler
	setRoutes()
	router.Use(nmiddleware.Log())
	router.Use(nmiddleware.Recover()) // in case of panic, recover and don't quit
	router.Use(middleware.RemoveTrailingSlash())
}

func setRoutes() {
	router.GET("/version", handler.Version)
}

// Run start the routeur
func Run(environment *env.Config) error {
	router.Server.Addr = environment.Address + ":" + strconv.Itoa(environment.Port)
	log.Infoln("Starting router on", router.Server.Addr)
	return graceful.ListenAndServe(router.Server, 10*time.Second)
}

// RunTLS start the routeur and use encryption to communicate
func RunTLS(environment *env.Config, certFile string, keyFile string) error {
	router.Server.Addr = environment.Address + ":" + strconv.Itoa(environment.Port)
	log.Infoln("Starting router on", router.Server.Addr)
	return graceful.ListenAndServeTLS(router.Server, certFile, keyFile, 10*time.Second)
}
