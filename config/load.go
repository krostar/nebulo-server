package config

import (
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/returncode"
)

// Load load the configuration
func Load() (err error) {
	Config = new(Options)
	parser = flags.NewParser(Config, flags.None)

	// load from command line
	if err = loadFromCommandLine(); err != nil {
		return err
	}

	// load default files
	for _, filepath := range defaultConfigurationFilePaths {
		if err = loadFromFile(filepath); err != nil {
			return fmt.Errorf("configuration file failed to load: %v", err)
		}
	}

	// try to override configuration via the command line
	if _, err = FromCommandLine(os.Args, Config); err != nil {
		return fmt.Errorf("configuration from command line failed to load: %v", err)
	}
	if err = loadSuccessfully("command line"); err != nil {
		return fmt.Errorf("configuration failed to load: %v", err)
	}

	if err = applyConfiguration(); err != nil {
		return fmt.Errorf("configuration application failed: %v", err)
	}
	return nil
}

func loadFromCommandLine() (err error) {
	// looking for help, config-gen or config-not-load-default and by command line syntax
	remainingArgs, err := FromCommandLine(os.Args, Config)
	if err != nil {
		return fmt.Errorf("Error while parsing configuration from command line: %s", err)
	}
	if len(remainingArgs) > 1 {
		return fmt.Errorf("Unknow argument: %s", strings.Join(remainingArgs[1:], " "))
	}
	if err = loadSuccessfully(""); err != nil {
		if err == errConfigHelp || err == errConfigGeneration {
			return nil
		}
		return fmt.Errorf("configuration failed to load: %v", err)
	}
	return nil
}

func loadFromFile(filename string) (err error) {
	if err = FromINIFile(filename, Config); err != nil {
		if os.IsNotExist(err) {
			log.Warningf("Error while parsing configuration from %s: does that file exist ?", filename)
		} else {
			return err
		}
	} else {
		if err = loadSuccessfully(filename); err != nil {
			if err == errConfigHelp || err == errConfigGeneration {
				return nil
			}
			return err
		}
	}
	return nil
}

func loadSuccessfully(from string) (err error) {
	if from != "" {
		log.Infof("Configuration successfully loaded from %s", from)
	}

	// if the user need help, the program has to quit
	if Config.Basic.Help {
		log.Infoln("Help parameter detected, program will print help on standart output and quit")
		parser = flags.NewParser(&Options{}, flags.None) // give a clear Options structure to avoid setted values interpreted as default values
		parser.WriteHelp(os.Stdout)
		os.Exit(returncode.HELP)
	}

	// this settings is used to have a template of all possible configuration, the program has to quit
	if Config.Basic.ConfigGeneration != "" {
		log.Infoln("Configuration generation parameter detected, program will generate configuration into default file and quit")
		if err := flags.NewIniParser(parser).WriteFile(Config.Basic.ConfigGeneration, flags.IniIncludeDefaults|flags.IniCommentDefaults|flags.IniIncludeComments); err != nil {
			return fmt.Errorf("flag parser creation failed: %v", err)
		}
		os.Exit(returncode.CONFIGGEN)
	}

	if Config.Configuration.DontLoadDefault {
		defaultConfigurationFilePaths = []string{}
	}

	// if we have the config parameter, we need to overload this config
	if confFile := Config.Configuration.File; confFile != "" {
		Config.Configuration.File = ""
		log.Warningln("The parameter -c or --config-file is detected, the specified file will override the current configuration")

		if err := FromINIFile(confFile, Config); err != nil {
			log.Warningln("Error while parsing configuration from specified configuration file "+confFile, err)
			return fmt.Errorf("unable to load configuration from file %s: %v", confFile, err)
		}
		return loadSuccessfully("specified configuration file " + confFile)
	}
	return nil
}
