package config

import (
	"errors"
	"fmt"

	gp "github.com/krostar/nebulo-golib/provider"
	gpMySQL "github.com/krostar/nebulo-golib/provider/mysql"
	gpSQLite "github.com/krostar/nebulo-golib/provider/sqlite"
	cp "github.com/krostar/nebulo-server/channel/provider"
	cpMySQL "github.com/krostar/nebulo-server/channel/provider/mysql"
	cpSQLite "github.com/krostar/nebulo-server/channel/provider/sqlite"
	mp "github.com/krostar/nebulo-server/message/provider"
	mpMySQL "github.com/krostar/nebulo-server/message/provider/mysql"
	mpSQLite "github.com/krostar/nebulo-server/message/provider/sqlite"
	up "github.com/krostar/nebulo-server/user/provider"
	upMySQL "github.com/krostar/nebulo-server/user/provider/mysql"
	upSQLite "github.com/krostar/nebulo-server/user/provider/sqlite"

	"github.com/krostar/nebulo-golib/log"
	"github.com/krostar/nebulo-server/env"
	validator "gopkg.in/validator.v2"
)

// Apply validate configuration and initialize needed package with
// values from configuration
func Apply() (err error) {
	// check the configuration
	if err = validator.Validate(Config); err != nil {
		return err
	}

	if err = ApplyLoggingOptions(&Config.Global.Logging); err != nil {
		return fmt.Errorf("apply logging configuration failed: %v", err)
	}

	applyEnvironmentOptions(&Config.Run.Environment)

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
	pdc := gp.DefaultConfig{
		CreateTablesIfNotExists: pc.CreateTablesIfNotExists,
		DropTablesIfExists:      pc.DropTablesIfExists,
	}

	switch pc.Type {
	case "sqlite":
		pc.SQLiteConfig.DefaultConfig = pdc
		err = gpSQLite.Use(&pc.SQLiteConfig)
	case "mysql":
		pc.MySQLConfig.DefaultConfig = pdc
		err = gpMySQL.Use(&pc.MySQLConfig)
	default:
		err = errors.New("unknown provider")
	}
	if err != nil {
		return fmt.Errorf("%s providers initialization failed: %v", pc.Type, err)
	}

	err = initProviders(pc)
	if err != nil {
		return fmt.Errorf("unable to initialized providers: %v", err)
	}

	err = resetProviders(&pdc)
	if err != nil {
		return fmt.Errorf("unable to reset providers: %v", err)
	}
	return nil
}

func initProviders(pc *providerOptions) (err error) {
	switch pc.Type {
	case "sqlite":
		if err = upSQLite.Init(); err != nil {
			return fmt.Errorf("sqlite user providers initialization failed: %v", err)
		}
		if err = cpSQLite.Init(); err != nil {
			return fmt.Errorf("sqlite channel providers initialization failed: %v", err)
		}
		if err = mpSQLite.Init(); err != nil {
			return fmt.Errorf("sqlite message providers initialization failed: %v", err)
		}
	case "mysql":
		if err = upMySQL.Init(); err != nil {
			return fmt.Errorf("mysql user providers initialization failed: %v", err)
		}
		if err = cpMySQL.Init(); err != nil {
			return fmt.Errorf("mysql channel providers initialization failed: %v", err)
		}
		if err = mpMySQL.Init(); err != nil {
			return fmt.Errorf("mysql message providers initialization failed: %v", err)
		}
	default:
		return fmt.Errorf("providers initialization failed: unknown %v provider", pc.Type)
	}

	log.Infof("users, channels and messages provided via %s", pc.Type)
	return nil
}

func resetProviders(pdc *gp.DefaultConfig) (err error) {
	if pdc.DropTablesIfExists {
		err = cp.P.DropTables()
		if err == nil {
			err = up.P.DropTables()
		}
		if err == nil {
			err = mp.P.DropTables()
		}
	}
	if err == nil && pdc.CreateTablesIfNotExists {
		err = cp.P.CreateTables()
		if err == nil {
			err = up.P.CreateTables()
		}
		if err == nil {
			err = mp.P.CreateTables()
		}
		if err == nil {
			err = cp.P.CreateIndexes()
		}
		if err == nil {
			err = up.P.CreateIndexes()
		}
		if err == nil {
			err = mp.P.CreateIndexes()
		}
	}
	return err
}
