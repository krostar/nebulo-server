package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	cli "gopkg.in/urfave/cli.v2"

	"github.com/krostar/nebulo-golib/log"
	"github.com/krostar/nebulo-server/config"
	"github.com/krostar/nebulo-server/router"
	"github.com/krostar/nebulo-server/router/handler"
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
}

func main() {
	app := &cli.App{
		Name:        "Nebulo",
		Usage:       "encrypted chat server",
		HideVersion: true,
		Before: func(c *cli.Context) (err error) {
			if err = config.ApplyLoggingOptions(&config.Config.Global.Logging); err != nil {
				return fmt.Errorf("unable to apply logging configuration: %v", err)
			}
			if configFile := c.String("config"); configFile != "" {
				if err = config.LoadFile(configFile); err != nil {
					return fmt.Errorf("unable to load configuration file %q:%v", configFile, err)
				}
			}
			return nil
		}, Flags: []cli.Flag{ // global flags (configuration and log purpose)
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to the configuration file",
			}, &cli.StringFlag{
				Name:        "log",
				Aliases:     []string{"l"},
				Usage:       "path to a file where the logs will be writted",
				DefaultText: "standart output",
				Destination: &config.CLI.Global.Logging.File,
			}, &cli.StringFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "level of informations to write (quiet, critical, error, warning, info, request, debug)",
				DefaultText: "debug",
				Destination: &config.CLI.Global.Logging.Verbose,
			},
		}, Commands: []*cli.Command{
			&cli.Command{ // run command, start the server
				Name:        "run",
				Usage:       "start the nebulo api server",
				Description: "required parameters description starts with a wildcard (*)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "environment",
						Aliases:     []string{"e"},
						Usage:       "* environment to use for external services connection purpose (dev, preprod, prod)",
						Destination: &config.CLI.Run.Environment.Type,
					}, &cli.StringFlag{
						Name:        "address",
						Aliases:     []string{"a"},
						Usage:       "override environment address to use to listen to",
						DefaultText: "depend on -e (environment)",
						Destination: &config.CLI.Run.Environment.Address,
					}, &cli.IntFlag{
						Name:        "port",
						Aliases:     []string{"p"},
						Usage:       "override environment port to use to listen to",
						DefaultText: "depend on -e (environment)",
						Destination: &config.CLI.Run.Environment.Port,
					}, &cli.StringFlag{
						Name:        "tls-crt",
						Usage:       "* tls certificate file used to encrypt communication (https)",
						Destination: &config.CLI.Run.TLS.Cert,
					}, &cli.StringFlag{
						Name:        "tls-key",
						Usage:       "* tls certificate key used with --tls-crt",
						Destination: &config.CLI.Run.TLS.Key,
					}, &cli.StringFlag{
						Name:        "tls-clients-ca",
						Usage:       "* tls certification authority used to validate clients certificate for the tls mutual authentication",
						Destination: &config.CLI.Run.TLS.ClientsCA.Cert,
					}, &cli.StringFlag{
						Name:        "tls-clients-ca-key",
						Usage:       "* tls certification authority key used with --tls-clients-ca",
						Destination: &config.CLI.Run.TLS.ClientsCA.Key,
					}, &cli.StringFlag{
						Name:        "tls-clients-ca-key-pwd",
						Usage:       "password/passphrase used with --tls-clients-ca-key",
						Destination: &config.CLI.Run.TLS.ClientsCA.KeyPassword,
					}, &cli.StringFlag{
						Name:        "provider",
						Usage:       "* database type to use to provide users and messages (sqlite)",
						Destination: &config.CLI.Run.Provider.Type,
					}, &cli.BoolFlag{
						Name:        "provider-createtables",
						Usage:       "create tables if not exists",
						DefaultText: "false",
						Destination: &config.CLI.Run.Provider.CreateTablesIfNotExists,
					}, &cli.BoolFlag{
						Name:        "provider-droptables",
						Usage:       "drop tables if not exists - only available in dev environment",
						DefaultText: "false",
						Destination: &config.CLI.Run.Provider.DropTablesIfExists,
					},
				}, Before: beforeCommandWhoNeedMergeConfiguration,
				Action: commandRun,
			}, &cli.Command{ // config-gen command, generate the configuration
				Name:  "config-gen",
				Usage: "generate a configuration file and quit",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "destination",
						Aliases:     []string{"d"},
						Usage:       "path to a file where the configuration will be writted",
						DefaultText: "standart output",
					},
				}, Before: beforeEveryCommand,
				Action: commandConfigGen,
			}, &cli.Command{ // version command output the version of the server
				Name:   "version",
				Usage:  "display the version",
				Before: beforeEveryCommand,
				Action: commandVersion,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Criticalf("unable to run app: %v", err)
		os.Exit(1)
	}
}

func beforeEveryCommand(c *cli.Context) (err error) {
	// we don't want useless arguments (non parsed arguments)
	if c.NArg() != 0 {
		return fmt.Errorf("unknown remaining args: %q", strings.Join(c.Args().Slice(), " "))
	}
	return nil
}

func beforeCommandWhoNeedMergeConfiguration(c *cli.Context) (err error) {
	if err = beforeEveryCommand(c); err != nil {
		return err
	}

	// merge configuration from cli and configuration file
	config.Merge()
	if err = config.Apply(); err != nil {
		return fmt.Errorf("configuration application failed: %v", err)
	}
	log.Logf(log.DEBUG, -1, "Configuration merged, validated and applied: %v", config.Config)
	return nil
}

func commandRun(_ *cli.Context) error {
	log.Infof("Starting Nebulo API server build %s (%s) on %s:%d", BuildVersion, BuildTime, config.Config.Run.Environment.Address, config.Config.Run.Environment.Port)
	return router.RunTLS(
		&config.Config.Run.Environment,
		config.Config.Run.TLS.Cert,
		config.Config.Run.TLS.Key,
		config.Config.Run.TLS.ClientsCA.Cert,
	)
}

func commandConfigGen(c *cli.Context) error {
	conf, err := json.MarshalIndent(config.Config, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to create json: %v", err)
	}
	if filepath := c.String("destination"); filepath != "" {
		if err := ioutil.WriteFile(filepath, conf, 0644); err != nil {
			return fmt.Errorf("unable to write sql queries file: %v", err)
		}
	} else {
		fmt.Println(string(conf))
	}
	return nil
}

func commandVersion(_ *cli.Context) error {
	fmt.Printf("nebulo %s (%s)\n", BuildVersion, BuildTime)
	return nil
}
