package main

import (
	"os"

	"github.com/krostar/nebulo/config"
	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/returncode"
	"github.com/krostar/nebulo/router"
)

func main() {

	log.Infoln("Starting nebulo api ....")

	if err := router.Run(env.Environment[config.Config.Environment]); err != nil {
		log.Criticalln(err)
		os.Exit(returncode.ROUTERFAILED)
	}

	log.Infoln("Server stop without errors")
	os.Exit(returncode.SUCCESS)
}
