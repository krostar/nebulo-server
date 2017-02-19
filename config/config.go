package config

import (
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/returncode"
	validator "gopkg.in/validator.v2"
)

//

// Options list all the available options of the program, with details useful for help command
type Options struct {
	Help        bool   `short:"h" long:"help" description:"show this help message" hidden:"true" validate:"-"`
	Config      string `short:"c" long:"config" description:"specify a configuration file (be cautious on infinite-recursive-configuration)"`
	Environment string `short:"e" long:"environment" choice:"dev" choice:"alpha" choice:"prod" description:"environment to use for external services connection purpose - this parameter is required" validate:"regexp=^(dev|alpha|prod)$"`
	Address     string `short:"a" long:"address" description:"override environment address to use to listen to" default-mask:"depend on -e (environment)" validate:"-"`
	Port        int    `short:"p" long:"port" description:"override port to use to listen to" default-mask:"depend on -e (environment)" validate:"-"`
	LogFile     string `short:"l" long:"logging-file" description:"the file where write the log (if not defined, logs are writted on standart output)" validate:"-"`
	Verbose     string `short:"v" long:"verbose" choice:"critical" choice:"error" choice:"warning" choice:"info" choice:"request" choice:"debug" description:"level of information write on standart output or in a file" default-mask:"debug"`
}

var (
	// Config is the active configuration of the program
	Config Options

	parser                 *flags.Parser
	atLeastOneConfigLoaded bool
)

func init() {
	parser = flags.NewParser(&Config, flags.None)
	var err error
	atLeastOneConfigLoaded = false

	// try to load configuration from default system location
	if err = FromINIFile("/etc/nebulo/config.ini", &Config); err != nil {
		log.Warningln("Error while parsing configuration from default directory config.ini: ", err)
	} else {
		successfullyLoaded(&Config, "default directory config.ini")
	}

	// try to override configuration via default folder location
	if err = FromINIFile("./config.ini", &Config); err != nil {
		log.Warningln("Error while parsing configuration from current directory config.ini: ", err)
	} else {
		successfullyLoaded(&Config, "current directory config.ini")
	}

	// try to override configuration via the command line
	if err = FromCommandLine(os.Args, &Config); err != nil {
		if atLeastOneConfigLoaded {
			log.Warningln("Error while parsing configuration from command line: ", err)
		} else {
			panic(err)
		}
	} else {
		successfullyLoaded(&Config, "command line")
	}

	// check the configuration
	if err = validator.Validate(Config); err != nil {
		log.Criticalln(err)
		panic(err)
	}

	// apply it
	if err = applyConfiguration(Config); err != nil {
		log.Criticalln(err)
		panic(err)
	}
}

func successfullyLoaded(config *Options, from string) {
	log.Infoln("Configuration successfully loaded from "+from, Config)

	// if the user need help, the program has to quit
	if config.Help {
		log.Infoln("Help parameter detected, program will print help on standart output and quit")
		parser = flags.NewParser(&Options{}, flags.None) // give a clear Options structure to avoid setted values interpreted as default values
		parser.WriteHelp(os.Stdout)
		os.Exit(returncode.HELP)
	}

	// if we have the config parameter, we need to overload this config
	if confFile := config.Config; confFile != "" {
		config.Config = ""
		log.Warningln("The parameter -c or --config is detected, the specified file will override the current configuration")

		if err := FromINIFile(confFile, config); err != nil {
			log.Errorln("Error while parsing configuration from specified configuration file "+confFile, err)
		} else {
			successfullyLoaded(config, "specified configuration file "+confFile)
			return
		}
	}
	atLeastOneConfigLoaded = true
}

func applyConfiguration(config Options) (err error) {
	// apply environment config
	if config.Environment != "" {
		environment := env.Environment[config.Environment]
		if config.Address != "" {
			environment.Address = config.Address
		}
		if config.Port > 0 {
			environment.Port = config.Port
		}
	}

	// apply log-related config
	if config.Verbose != "" {
		log.Verbosity = log.VerboseMapping[config.Verbose]
	}
	if config.LogFile != "" {
		if err = log.SetFile(config.LogFile); err != nil {
			return err
		}
	}

	return nil
}

// FromCommandLine load options from program arguments
func FromCommandLine(args []string, conf *Options) (err error) {
	_, err = parser.ParseArgs(args)
	return err
}

// FromINIFile load options from configuration file
func FromINIFile(filename string, conf *Options) (err error) {
	return flags.IniParse(filename, conf)
}
