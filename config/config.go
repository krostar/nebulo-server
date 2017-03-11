package config

import (
	"errors"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/krostar/nebulo/tools/validator" // used to init custom validators before using them
)

// Options list all the available options of the program, with details useful for help command and validators to help validations of fields
type Options struct {
	ConfigGeneration string `long:"config-gen" description:"generate a configuration file for the actual configuration to the specified file and quit" no-ini:"true" validate:"-"`
	Help             bool   `short:"h" long:"help" description:"show this help message" no-ini:"true" validate:"-"`

	ConfigDontLoadDefault bool   `long:"config-dont-load-default" description:"choose to load or not the default configuration files" no-ini:"true" validate:"-"`
	ConfigFile            string `short:"c" long:"config-file" description:"specify a configuration file (be cautious on infinite-recursive-configuration)" validate:"file=omitempty+readable"`

	Environment string `short:"e" long:"environment" choice:"dev" choice:"beta" choice:"prod" description:"environment to use for external services connection purpose - this parameter is required" validate:"regexp=^(dev|beta|prod)$"`
	Address     string `short:"a" long:"address" description:"override environment address to use to listen to" default-mask:"depend on -e (environment)"`
	Port        int    `short:"p" long:"port" description:"override environment port to use to listen to" default-mask:"depend on -e (environment)"`

	TLSCertFile             string `long:"tls-crt-file" description:"tls certificate file used to encrypt communication - this parameter is required for TLS communication" validate:"file=readable"`
	TLSKeyFile              string `long:"tls-key-file" description:"tls certificate key used to encrypt communication - this parameter is required for TLS communication" validate:"file=readable"`
	TLSClientsCACertFile    string `long:"tls-clients-ca-cert-file" description:"tls certification authority used to validate clients certificate for the tls mutual authentication - this parameter is required for TLS communication" validate:"file=readable"`
	TLSClientsCAKeyFile     string `long:"tls-clients-ca-key-file" description:"tls certification authority key used to validate clients certificate for the tls mutual authentication - this parameter is required for TLS communication" validate:"file=readable"`
	TLSClientsCAKeyPassword string `long:"tls-clients-ca-key-pwd" description:"tls certification authority key password used to validate clients certificate for the tls mutual authentication - this parameter is required for TLS communication"`

	LogFile string `short:"l" long:"logging-file" description:"the file where write the log" default-mask:"no file, standart output" validate:"-"`
	Verbose string `short:"v" long:"verbose" choice:"quiet" choice:"critical" choice:"error" choice:"warning" choice:"info" choice:"request" choice:"debug" description:"level of information to write on standart output or in a file" default-mask:"debug" validate:"regexp=^(quiet|critical|error|warning|info|request|debug)?$"`

	UserProvider     string `long:"user-provider" choice:"file" description:"provider to use to get users informations" validate:"regexp=^(file)?$"`
	UserProviderFile string `long:"user-provider-file" description:"provider file path where users informations are stored" validate:"-"`
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

// FromCommandLine load options from program arguments
func FromCommandLine(args []string, conf *Options) (remaining []string, err error) {
	return parser.ParseArgs(args)
}

// FromINIFile load options from configuration file
func FromINIFile(filename string, conf *Options) (err error) {
	return flags.IniParse(filename, conf)
}
