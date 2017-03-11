package config

import (
	"errors"
	"fmt"

	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
	up "github.com/krostar/nebulo/user/provider"
	upf "github.com/krostar/nebulo/user/provider/file"
	validator "gopkg.in/validator.v2"
)

func applyConfiguration() (err error) {
	// check the configuration
	if err = validator.Validate(Config); err != nil {
		return err
	}

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
			return fmt.Errorf("unable to set log outputfile: %v", err)
		}
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

func applyProviderConfiguration() (p up.Provider, err error) {
	switch Config.UserProvider {
	case "file":
		p, err = upf.NewFromConfig(&upf.Config{
			Filepath: Config.UserProviderFile,
		})
		if err != nil {
			return nil, fmt.Errorf("user-provider-file creation failed: %v", err)
		}
	default:
		return nil, errors.New("unknown user provider")
	}

	log.Debugf("Using %s to provide user", Config.UserProvider)
	return p, nil
}
