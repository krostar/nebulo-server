package config

import (
	"errors"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/krostar/nebulo/tools/validator" // used to init custom validators before using them
)

type basicOptions struct {
	Help             bool `short:"h" long:"help" description:"show this help message" no-ini:"true" validate:"-"`
	ConfigGeneration bool `long:"config-gen" description:"generate a configuration for the actual parameter, print on standart output and quit" no-ini:"true" validate:"-"`
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
	CertFile  string       `long:"tls-crt-file" description:"tls certificate file used to encrypt communication - this parameter is required" validate:"file=readable"`
	KeyFile   string       `long:"tls-key-file" description:"tls certificate key used to encrypt communication - this parameter is required" validate:"file=readable"`
	ClientsCA tlsClientsCA `group:"Clients CA"`
}

type tlsClientsCA struct {
	CertFile    string `long:"tls-clients-ca-cert-file" description:"tls certification authority used to validate clients certificate for the tls mutual authentication - this parameter is required" validate:"file=readable"`
	KeyFile     string `long:"tls-clients-ca-key-file" description:"tls certification authority key file used to validate clients certificate for the tls mutual authentication - this parameter is required" validate:"file=readable"`
	KeyPassword string `long:"tls-clients-ca-key-pwd" description:"tls certification authority key password used to validate clients certificate for the tls mutual authentication" default-mask:"no password"`
}

type sqlOptions struct {
	Type string `long:"sql-type" choice:"sqlite" description:"provider to use to get users informations - this parameter is required" validate:"regexp=^(sqlite)?$"`

	CreateTablesIfNotExists bool `long:"sql-createtables" description:"create tables if not exists" default-mask:"false" validate:"-"`
	DropTablesIfExists      bool `long:"sql-droptables" description:"drop tables if exists" default-mask:"false" validate:"-"`
	CreateQuery             bool `long:"sql-createqueries" description:"print the sql create queries and quit" no-ini:"true" validate:"-"`

	SQLiteFile string `long:"sql-sqlite-file" description:"sqlite filepath where informations are stored" validate:"-"`
}

// Options list all the available options of the program, with details useful for help command and validators to help validations of fields
type Options struct {
	Basic         basicOptions
	Configuration configurationOptions `group:"Configuration Options"`
	Environment   environmentOptions   `group:"Environment Options"`
	Logging       loggingOptions       `group:"Logging Options"`
	TLS           tlsOptions           `group:"TLS Options"`
	SQL           sqlOptions           `group:"SQL Options"`
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
