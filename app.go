package main

import (
	"os"

	"github.com/krostar/nebulo/config"
	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/handler"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/returncode"
	"github.com/krostar/nebulo/router"
)

var (
	// BuildTime represent the time when the binary has been created
	BuildTime = "undefined"
	// BuildVersion represent the version of the binary (git tag or revision)
	BuildVersion = "undefined"
)

func init() {
	handler.BuildTime = BuildTime
	handler.BuildVersion = BuildVersion
}

func main() {

	log.Infoln("Starting nebulo api", BuildVersion)

	var err error
	if config.Config.TLSCertFile != "" && config.Config.TLSKeyFile != "" {
		err = router.RunTLS(env.Environment[config.Config.Environment], config.Config.TLSCertFile, config.Config.TLSKeyFile)
	} else {
		err = router.Run(env.Environment[config.Config.Environment])
	}
	if err != nil {
		log.Criticalln(err)
		os.Exit(returncode.ROUTERFAILED)
	}

	log.Infoln("Server stop without errors")
	os.Exit(returncode.SUCCESS)
}
