package config

import (
	"errors"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/krostar/nebulo/tools/validator" // used to init custom validators before using them
)

type basicOptions struct {
	Help             bool   `short:"h" long:"help" description:"show this help message" no-ini:"true" validate:"-"`
	ConfigGeneration string `long:"config-gen" description:"generate a configuration file for the actual configuration to the specified file and quit" no-ini:"true" validate:"-"`
}

type configurationOptions struct {
	DontLoadDefault bool   `long:"config-dont-load-default" description:"choose to load or not the default configuration files" no-ini:"true" default-mask:"false" validate:"-"`
	File            string `short:"c" long:"config-file" description:"specify a configuration file (be cautious on infinite-recursive-configuration)" no-ini:"true" validate:"file=omitempty+readable"`
}

type environmentOptions struct {
	Type    string `short:"e" long:"environment" choice:"dev" choice:"beta" choice:"prod" description:"environment to use for external services connection purpose - this parameter is required" validate:"regexp=^(dev|beta|prod)$"`
	Address string `short:"a" long:"address" description:"override environment address to use to listen to" default-mask:"depend on -e (environment)"`
	Port    int    `short:"p" long:"port" description:"override environment port to use to listen to" default-mask:"depend on -e (environment)"`
}

type loggingOptions struct {
	File    string `short:"l" long:"logging-file" description:"the file where write the log" default-mask:"no file, standart output" validate:"-"`
	Verbose string `short:"v" long:"verbose" choice:"quiet" choice:"critical" choice:"error" choice:"warning" choice:"info" choice:"request" choice:"debug" description:"level of information to write on standart output or in a file" default-mask:"debug" validate:"regexp=^(quiet|critical|error|warning|info|request|debug)?$"`
}

type tlsOptions struct {
	CertFile             string `long:"tls-crt-file" description:"tls certificate file used to encrypt communication - this parameter is required" validate:"file=readable"`
	KeyFile              string `long:"tls-key-file" description:"tls certificate key used to encrypt communication - this parameter is required" validate:"file=readable"`
	ClientsCACertFile    string `long:"tls-clients-ca-cert-file" description:"tls certification authority used to validate clients certificate for the tls mutual authentication - this parameter is required" validate:"file=readable"`
	ClientsCAKeyFile     string `long:"tls-clients-ca-key-file" description:"tls certification authority key file used to validate clients certificate for the tls mutual authentication - this parameter is required" validate:"file=readable"`
	ClientsCAKeyPassword string `long:"tls-clients-ca-key-pwd" description:"tls certification authority key password used to validate clients certificate for the tls mutual authentication" default-mask:"no password"`
}

type userProviderOptions struct {
	Type                    string `long:"user-provider" choice:"sqlite" description:"provider to use to get users informations" validate:"regexp=^(sqlite)?$"`
	CreateTablesIfNotExists bool   `long:"user-provider-createtable" description:"create tables if not exists" default-mask:"false" validate:"-"`
	DropTablesIfExists      bool   `long:"user-provider-droptables" description:"drop tables if exists" default-mask:"false" validate:"-"`
	SQLiteFile              string `long:"user-provider-sqlite-file" description:"provider sqlite filepath where users informations are stored" validate:"-"`
}

// Options list all the available options of the program, with details useful for help command and validators to help validations of fields
type Options struct {
	Basic         basicOptions
	Configuration configurationOptions `group:"Configuration Options"`
	Environment   environmentOptions   `group:"Environment Options"`
	Logging       loggingOptions       `group:"Logging Options"`
	TLS           tlsOptions           `group:"TLS Options"`
	UserProvider  userProviderOptions  `group:"User Provider Options"`
}

var (
	// Config is the active configuration of the program
	Config *Options

	parser                        *flags.Parser
	defaultConfigurationFilePaths []string
	errConfigHelp                 = errors.New("config help")
	errConfigGeneration           = errors.New("config generation")
)

func init() {
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
