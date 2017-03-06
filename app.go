package main

import (
	"os"

	"github.com/krostar/nebulo/config"
	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/returncode"
	"github.com/krostar/nebulo/router"
	"github.com/krostar/nebulo/router/handler"
)

var (
	// BuildTime represent the time when the binary has been created
	BuildTime = "undefined"
	// BuildVersion is the version of the binary (git tag or revision)
	BuildVersion = "undefined"
)

func init() {
	handler.BuildTime = BuildTime
	handler.BuildVersion = BuildVersion

	if err := config.Load(); err != nil {
		log.Criticalf("Unable to load configuration: %v", err)
		os.Exit(returncode.CONFIGFAILED)
	}
}

func main() {
	log.Infof("Starting nebulo api %q", BuildVersion)

	if err := router.RunTLS(
		env.EnvironmentConfig[env.Environment(config.Config.Environment)],
		config.Config.TLSCertFile,
		config.Config.TLSKeyFile,
		config.Config.TLSClientsCACertFile,
	); err != nil {
		log.Criticalln(err)
		os.Exit(returncode.ROUTERFAILED)
	}

	log.Infoln("Server stop without errors")
	os.Exit(returncode.SUCCESS)
}
