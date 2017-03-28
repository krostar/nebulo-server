package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/returncode"
	up "github.com/krostar/nebulo/user/provider"
	upLite "github.com/krostar/nebulo/user/provider/sqlite"
	validator "gopkg.in/validator.v2"
)

func applyConfiguration() (err error) {
	// check the configuration
	if err = validator.Validate(Config); err != nil {
		return err
	}

	// apply environment config
	err = applyEnvironmentConfiguration()
	if err != nil {
		return fmt.Errorf("apply provider configuration failed: %v", err)
	}

	// apply log-related config
	if err = applyLoggingConfiguration(); err != nil {
		return fmt.Errorf("apply logging configuration failed: %v", err)
	}

	// apply provider-related config
	err = applyProviderConfiguration()
	if err != nil {
		return fmt.Errorf("apply provider configuration failed: %v", err)
	}

	return nil
}

func applyEnvironmentConfiguration() (err error) {
	if Config.Environment.Type != "" {
		environment := env.EnvironmentConfig[env.Environment(Config.Environment.Type)]
		if Config.Environment.Address != "" {
			environment.Address = Config.Environment.Address
		}
		if Config.Environment.Port > 0 {
			environment.Port = Config.Environment.Port
		}
		return nil
	}
	return errors.New("unknown environment")
}

func applyLoggingConfiguration() (err error) {
	if Config.Logging.Verbose != "" {
		log.Verbosity = log.VerboseMapping[Config.Logging.Verbose]
	}
	if Config.Logging.File != "" {
		if err = log.SetOutputFile(Config.Logging.File); err != nil {
			return fmt.Errorf("unable to set log outputfile: %v", err)
		}
	}
	return nil
}

func applyProviderConfiguration() (err error) {
	var uP up.Provider

	switch Config.UserProvider.Type {
	case "sqlite":
		uP, err = upLite.NewFromConfig(&upLite.Config{
			Filepath:                Config.UserProvider.SQLiteFile,
			CreateTablesIfNotExists: Config.UserProvider.CreateTablesIfNotExists,
			DropTablesIfExists:      Config.UserProvider.DropTablesIfExists,
		})
		if err != nil {
			return fmt.Errorf("user-provider-sqlite initialization failed: %v", err)
		}
	default:
		return errors.New("unknown user provider")
	}

	if err = up.Use(uP); err != nil {
		return fmt.Errorf("unable to set user provider: %v", err)
	}
	log.Infof("Using %s to provide user", Config.UserProvider.Type)

	if Config.UserProvider.SQLCreateQuery {
		log.Infoln("User SQL creation query parameter detected, program will output query and quit")
		creationQuery, err := up.P.SQLCreateQuery()
		if err != nil {
			return fmt.Errorf("unable to get sql user provider creation query: %v", err)
		}
		fmt.Println(creationQuery)
		os.Exit(returncode.SQLGEN)
	}

	return nil
}
