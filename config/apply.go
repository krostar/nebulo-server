package config

import (
	"errors"
	"fmt"

	cp "github.com/krostar/nebulo/channel/provider"
	cpSQLite "github.com/krostar/nebulo/channel/provider/sqlite"
	"github.com/krostar/nebulo/env"
	"github.com/krostar/nebulo/log"
	gp "github.com/krostar/nebulo/provider"
	up "github.com/krostar/nebulo/user/provider"
	upSQLite "github.com/krostar/nebulo/user/provider/sqlite"
	validator "gopkg.in/validator.v2"
)

// Apply validate configuration and initialize needed package with
// values from configuration
func Apply() (err error) {
	// check the configuration
	if err = validator.Validate(Config); err != nil {
		return err
	}

	// apply log-related config
	if err = ApplyLoggingOptions(&Config.Global.Logging); err != nil {
		return fmt.Errorf("apply logging configuration failed: %v", err)
	}

	applyEnvironmentOptions(&Config.Run.Environment)

	// apply provider-related config
	err = ApplyProvidersOptions(&Config.Run.Provider)
	if err != nil {
		return fmt.Errorf("apply providers configuration failed: %v", err)
	}

	return nil
}

// ApplyLoggingOptions apply configuration on log package
func ApplyLoggingOptions(lc *logOptions) (err error) {
	if lc.Verbose != "" {
		log.Verbosity = log.VerboseMapping[lc.Verbose]
	}
	if lc.File != "" {
		if err = log.SetOutputFile(lc.File); err != nil {
			return fmt.Errorf("unable to set log outputfile: %v", err)
		}
	}
	return nil
}

func applyEnvironmentOptions(ec *env.Config) {
	var (
		addr = env.EnvironmentConfig[ec.Type].Address
		port = env.EnvironmentConfig[ec.Type].Port
	)

	if ec.Address == "" {
		ec.Address = addr
	}
	if ec.Port == 0 {
		ec.Port = port
	}
}

// ApplyProvidersOptions apply configuration on providers package
func ApplyProvidersOptions(pc *providerOptions) (err error) {
	var (
		uP up.Provider
		cP cp.Provider
	)

	switch pc.Type {
	case "sqlite":
		// init users provider
		uP, err = upSQLite.NewFromConfig(&gp.SQLiteConfig{
			Filepath: pc.SQLiteFile,
			DefaultConfig: gp.DefaultConfig{
				CreateTablesIfNotExists: pc.CreateTablesIfNotExists,
				DropTablesIfExists:      pc.DropTablesIfExists,
			},
		})
		if err != nil {
			return fmt.Errorf("users provider sqlite initialization failed: %v", err)
		}

		// init channels provider
		cP, err = cpSQLite.NewFromConfig(&gp.SQLiteConfig{
			Filepath: pc.SQLiteFile,
			DefaultConfig: gp.DefaultConfig{
				CreateTablesIfNotExists: pc.CreateTablesIfNotExists,
				DropTablesIfExists:      pc.DropTablesIfExists,
			},
		})
		if err != nil {
			return fmt.Errorf("channels provider sqlite initialization failed: %v", err)
		}

	default:
		return errors.New("unknown provider")
	}

	if err = up.Use(uP); err != nil {
		return fmt.Errorf("unable to set users provider: %v", err)
	}
	log.Infof("Using %s to provide users", pc.Type)

	if err = cp.Use(cP); err != nil {
		return fmt.Errorf("unable to set channels provider: %v", err)
	}
	log.Infof("Using %s to provide channels", pc.Type)

	return nil
}
