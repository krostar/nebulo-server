package config

import (
	"errors"
	"fmt"
	"os"

	cp "github.com/krostar/nebulo/channel/provider"
	cpSQLite "github.com/krostar/nebulo/channel/provider/sqlite"
	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
	gp "github.com/krostar/nebulo/provider"
	"github.com/krostar/nebulo/returncode"
	up "github.com/krostar/nebulo/user/provider"
	upSQLite "github.com/krostar/nebulo/user/provider/sqlite"
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
	err = applyProvidersConfiguration()
	if err != nil {
		return fmt.Errorf("apply providers configuration failed: %v", err)
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

func applyProvidersConfiguration() (err error) {
	if err = applyUsersProviderConfiguration(); err != nil {
		return fmt.Errorf("unable to apply users provider configuration: %v", err)
	}
	if err = applyChannelsProviderConfiguration(); err != nil {
		return fmt.Errorf("unable to apply channels provider configuration: %v", err)
	}

	if Config.SQL.CreateQuery {
		log.Infoln("SQL creation query parameter detected, program will output queries and quit")
		userCreationQuery, err := up.P.SQLCreateQuery()
		if err != nil {
			return fmt.Errorf("unable to get sql users provider creation query: %v", err)
		}
		channelCreationQuery, err := cp.P.SQLCreateQuery()
		if err != nil {
			return fmt.Errorf("unable to get sql channels provider creation query: %v", err)
		}
		fmt.Printf("%s\n\n%s\n", userCreationQuery, channelCreationQuery)
		os.Exit(returncode.SQLGEN)
	}

	return nil
}

func applyUsersProviderConfiguration() (err error) {
	var uP up.Provider

	switch Config.SQL.Type {
	case "sqlite":
		// init users provider
		uP, err = upSQLite.NewFromConfig(&gp.SQLiteConfig{
			Filepath:                Config.SQL.SQLiteFile,
			CreateTablesIfNotExists: Config.SQL.CreateTablesIfNotExists,
			DropTablesIfExists:      Config.SQL.DropTablesIfExists,
		})
		if err != nil {
			return fmt.Errorf("users provider sqlite initialization failed: %v", err)
		}
	default:
		return errors.New("unknown provider")
	}

	if err = up.Use(uP); err != nil {
		return fmt.Errorf("unable to set users provider: %v", err)
	}
	log.Infof("Using %s to provide users", Config.SQL.Type)

	return nil
}

func applyChannelsProviderConfiguration() (err error) {
	var cP cp.Provider

	switch Config.SQL.Type {
	case "sqlite":
		// init users provider
		cP, err = cpSQLite.NewFromConfig(&gp.SQLiteConfig{
			Filepath:                Config.SQL.SQLiteFile,
			CreateTablesIfNotExists: Config.SQL.CreateTablesIfNotExists,
			DropTablesIfExists:      Config.SQL.DropTablesIfExists,
		})
		if err != nil {
			return fmt.Errorf("channels provider sqlite initialization failed: %v", err)
		}
	default:
		return errors.New("unknown provider")
	}

	if err = cp.Use(cP); err != nil {
		return fmt.Errorf("unable to set channels provider: %v", err)
	}
	log.Infof("Using %s to provide channels", Config.SQL.Type)

	return nil
}
