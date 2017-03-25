package config

import (
	"errors"
	"fmt"

	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
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
	_, err = applyEnvironmentConfiguration()
	if err != nil {
		return fmt.Errorf("apply provider configuration failed: %v", err)
	}

	// apply log-related config
	err = applyLoggingConfiguration()
	if err != nil {
		return fmt.Errorf("apply logging configuration failed: %v", err)
	}

	// apply provider-related config
	userProvider, err := applyProviderConfiguration()
	if err != nil {
		return fmt.Errorf("apply provider configuration failed: %v", err)
	}
	if err = up.Use(userProvider); err != nil {
		return fmt.Errorf("unable to set user provider: %v", err)
	}

	return nil
}

func applyEnvironmentConfiguration() (environment *env.Config, err error) {
	if Config.Environment.Type != "" {
		environment = env.EnvironmentConfig[env.Environment(Config.Environment.Type)]
		if Config.Environment.Address != "" {
			environment.Address = Config.Environment.Address
		}
		if Config.Environment.Port > 0 {
			environment.Port = Config.Environment.Port
		}
		return environment, nil
	}
	return nil, errors.New("unknown environment")
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

func applyProviderConfiguration() (p up.Provider, err error) {
	switch Config.UserProvider.Type {
	case "sqlite":
		p, err = upLite.NewFromConfig(&upLite.Config{
			Filepath:                Config.UserProvider.SQLiteFile,
			CreateTablesIfNotExists: Config.UserProvider.CreateTablesIfNotExists,
			DropTablesIfExists:      Config.UserProvider.DropTablesIfExists,
		})
		if err != nil {
			return nil, fmt.Errorf("user-provider-file creation failed: %v", err)
		}
	default:
		return nil, errors.New("unknown user provider")
	}

	log.Debugf("Using %s to provide user", Config.UserProvider.Type)
	return p, nil
}
