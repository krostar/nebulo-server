package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/tools"
	_ "github.com/krostar/nebulo/tools/validator" // used to init custom validators before using them
)

type globalOptions struct {
	Logging logOptions `json:"log"`
}

type logOptions struct {
	Verbose string `json:"verbose" validate:"regexp=^(quiet|critical|error|warning|info|request|debug)?$"`
	File    string `json:"file" validate:"file=omitempty+readable"`
}

type runOptions struct {
	Environment env.Config      `json:"env"`
	TLS         tlsOptions      `json:"tls"`
	Provider    providerOptions `json:"provider"`
}

type tlsOptions struct {
	Cert      string       `json:"cert" validate:"file=readable"`
	Key       string       `json:"key" validate:"file=readable"`
	ClientsCA tlsClientsCA `json:"clients-ca" validate:"-"`
}

type tlsClientsCA struct {
	Cert        string `json:"cert" validate:"file=readable"`
	Key         string `json:"key" validate:"file=readable"`
	KeyPassword string `json:"key-password" validate:"-"`
}

type providerOptions struct {
	Type string `json:"type" validate:"regexp=^(sqlite)?$"`

	CreateTablesIfNotExists bool `json:"-" validate:"-"`
	TruncateTables          bool `json:"-" validate:"-"`
	DropTablesIfExists      bool `json:"-" validate:"-"`

	SQLiteFile string `json:"sqlite-file" validate:"file=omitempty+readable"`
}

// Options list all the available configurations
type Options struct {
	Global globalOptions `json:"global"`
	Run    runOptions    `json:"run"`
}

var (
	// Config store the active merge configuration
	Config = &Options{}
	// CLI store the configuration fetched from the console line
	CLI = &Options{}
	// File store the configuration fetched from an optional file
	File = &Options{}
)

// LoadFile fill config.File with the configuration parsed from path
func LoadFile(path string) (err error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}

	if err = json.Unmarshal(raw, File); err != nil {
		return fmt.Errorf("unable to parse json file: %v", err)
	}

	return nil
}

// Merge fill config.Config based on config.CLI and config.File
// File < CLI
func Merge() {
	mergeRecursive(reflect.ValueOf(CLI).Elem(), reflect.ValueOf(File).Elem(), reflect.ValueOf(Config).Elem())
}

func mergeRecursive(cli, file, config reflect.Value) {
	switch config.Kind() {
	case reflect.Struct:
		for i := 0; i < config.NumField(); i++ {
			mergeRecursive(cli.Field(i), file.Field(i), config.Field(i))
		}
	default:
		if !tools.IsZeroOrNil(cli) && cli.String() != "" {
			config.Set(cli)
		} else if !tools.IsZeroOrNil(file) && file.String() != "" {
			config.Set(file)
		}
	}
}
