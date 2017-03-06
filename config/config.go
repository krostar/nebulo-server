package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/returncode"
	_ "github.com/krostar/nebulo/tools/validator" // used to init custom validators before using them
	validator "gopkg.in/validator.v2"
)

// Options list all the available options of the program, with details useful for help command and validators to help validations of fields
type Options struct {
	Help                    bool   `short:"h" long:"help" description:"show this help message" no-ini:"true" validate:"-"`
	ConfigDontLoadDefault   bool   `long:"config-dont-load-default" description:"choose to load or not the default configuration files" no-ini:"true" validate:"-"`
	ConfigGeneration        string `long:"config-gen" description:"generate a configuration file for the actual configuration to the specified file and quit" no-ini:"true" validate:"-"`
	ConfigFile              string `short:"c" long:"config-file" description:"specify a configuration file (be cautious on infinite-recursive-configuration)" validate:"file=omitempty+readable"`
	Environment             string `short:"e" long:"environment" choice:"dev" choice:"beta" choice:"prod" description:"environment to use for external services connection purpose - this parameter is required" validate:"regexp=^(dev|beta|prod)$"`
	Address                 string `short:"a" long:"address" description:"override environment address to use to listen to" default-mask:"depend on -e (environment)"`
	Port                    int    `short:"p" long:"port" description:"override environment port to use to listen to" default-mask:"depend on -e (environment)"`
	TLSCertFile             string `long:"tls-crt-file" description:"tls certificate file used to encrypt communication - this parameter is required for TLS communication" validate:"file=readable"`
	TLSKeyFile              string `long:"tls-key-file" description:"tls certificate key used to encrypt communication - this parameter is required for TLS communication" validate:"file=readable"`
	TLSClientsCACertFile    string `long:"tls-clients-ca-cert-file" description:"tls certification authority used to validate clients certificate for the tls mutual authentication - this parameter is required for TLS communication" validate:"file=readable"`
	TLSClientsCAKeyFile     string `long:"tls-clients-ca-key-file" description:"tls certification authority key used to validate clients certificate for the tls mutual authentication - this parameter is required for TLS communication" validate:"file=readable"`
	TLSClientsCAKeyPassword string `long:"tls-clients-ca-key-pwd" description:"tls certification authority key password used to validate clients certificate for the tls mutual authentication - this parameter is required for TLS communication"`
	LogFile                 string `short:"l" long:"logging-file" description:"the file where write the log" default-mask:"no file, standart output" validate:"-"`
	Verbose                 string `short:"v" long:"verbose" choice:"quiet" choice:"critical" choice:"error" choice:"warning" choice:"info" choice:"request" choice:"debug" description:"level of information to write on standart output or in a file" default-mask:"debug" validate:"regexp=^(quiet|critical|error|warning|info|request|debug)?$"`
}

var (
	// Config is the active configuration of the program
	Config *Options

	parser                        *flags.Parser
	tester                        bool
	defaultConfigurationFilePaths []string

	errConfigHelp       = errors.New("config help")
	errConfigGeneration = errors.New("config generation")
)

func init() {
	tester = false
	defaultConfigurationFilePaths = []string{
		"/etc/nebulo/config.ini",
		"./config.ini",
	}
}

// Load load the configuration
func Load() (err error) {
	Config = new(Options)
	parser = flags.NewParser(Config, flags.None)

	// looking for help, config-gen or config-not-load-default and by command line syntax
	remainingArgs, err := FromCommandLine(os.Args, Config)
	if err != nil {
		return fmt.Errorf("Error while parsing configuration from command line: %s", err)
	}
	if len(remainingArgs) > 1 {
		return fmt.Errorf("Unknow argument: %s", strings.Join(remainingArgs[1:], " "))
	}
	if err = successfullyLoaded(""); err != nil {
		if err == errConfigHelp || err == errConfigGeneration {
			return nil
		}
		return err
	}

	// load default files
	for _, filepath := range defaultConfigurationFilePaths {
		if err = fromFile(filepath); err != nil {
			return err
		}
	}

	// try to override configuration via the command line
	if _, err = FromCommandLine(os.Args, Config); err != nil {
		return err
	}
	if err = successfullyLoaded("command line"); err != nil {
		return err
	}

	if err = applyConfiguration(); err != nil {
		return err
	}
	return nil
}

func successfullyLoaded(from string) (err error) {
	if from != "" {
		log.Infof("Configuration successfully loaded from %s", from)
	}

	// if the user need help, the program has to quit
	if Config.Help {
		log.Infoln("Help parameter detected, program will print help on standart output and quit")
		parser = flags.NewParser(&Options{}, flags.None) // give a clear Options structure to avoid setted values interpreted as default values
		parser.WriteHelp(os.Stdout)
		if !tester {
			os.Exit(returncode.HELP)
		} else {
			return errConfigHelp
		}
	}

	// this settings is used to have a template of all possible configuration, the program has to quit
	if Config.ConfigGeneration != "" {
		log.Infoln("Configuration generation parameter detected, program will generate configuration into default file and quit")
		if err := flags.NewIniParser(parser).WriteFile(Config.ConfigGeneration, flags.IniIncludeDefaults|flags.IniCommentDefaults|flags.IniIncludeComments); err != nil {
			return err
		}
		if !tester {
			os.Exit(returncode.CONFIGGEN)
		} else {
			return errConfigGeneration
		}
	}

	if Config.ConfigDontLoadDefault {
		defaultConfigurationFilePaths = []string{}
	}

	// if we have the config parameter, we need to overload this config
	if confFile := Config.ConfigFile; confFile != "" {
		Config.ConfigFile = ""
		log.Warningln("The parameter -c or --config-file is detected, the specified file will override the current configuration")

		if err := FromINIFile(confFile, Config); err != nil {
			log.Warningln("Error while parsing configuration from specified configuration file "+confFile, err)
			return err
		}
		return successfullyLoaded("specified configuration file " + confFile)
	}
	return nil
}

func applyConfiguration() (err error) {
	// check the configuration
	if err = validator.Validate(Config); err != nil {
		return err
	}

	// check TLS configuration
	// if Config.TLSCertFile != "" && (Config.TLSKeyFile == "" || Config.TLSClientsCAFile == "") {
	// 	return errors.New("TLSCertFile, TLSKeyFile and TLSClientsCA are required for TLS communication")
	// }

	// apply environment config
	if Config.Environment != "" {
		environment := env.EnvironmentConfig[env.Environment(Config.Environment)]
		if Config.Address != "" {
			environment.Address = Config.Address
		}
		if Config.Port > 0 {
			environment.Port = Config.Port
		}
	}

	// apply log-related config
	if Config.Verbose != "" {
		log.Verbosity = log.VerboseMapping[Config.Verbose]
	}
	if Config.LogFile != "" {
		if err = log.SetOutputFile(Config.LogFile); err != nil {
			return err
		}
	}

	return nil
}

func fromFile(filename string) (err error) {
	if err = FromINIFile(filename, Config); err != nil {
		if os.IsNotExist(err) {
			log.Warningf("Error while parsing configuration from %s: does that file exist ?", filename)
		} else {
			return fmt.Errorf("unable to parse file: %s", err)
		}
	} else {
		if err = successfullyLoaded(filename); err != nil {
			if err == errConfigHelp || err == errConfigGeneration {
				return nil
			}
			return err
		}
	}
	return nil
}

// FromCommandLine load options from program arguments
func FromCommandLine(args []string, conf *Options) (remaining []string, err error) {
	return parser.ParseArgs(args)
}

// FromINIFile load options from configuration file
func FromINIFile(filename string, conf *Options) (err error) {
	return flags.IniParse(filename, conf)
}
