package main

import (
	"flag"
	"time"

	"github.com/Hentra/dyndns-netcup-go/internal"
)

const (
	defaultConfigFile = "config.yml"
	configUsage       = "Specify location of the config file"
	verboseUsage      = "Use verbose output"
)

type cmdConfig struct {
	ConfigFile string
	Verbose    bool
}

func main() {
	cmdConfig := parseCmd()

	logger := internal.NewLogger(cmdConfig.Verbose)

	config, err := internal.LoadConfig(cmdConfig.ConfigFile)
	if err != nil {
		logger.Error(err)
	}

	cache, err := internal.NewCache(config.IPCache, time.Second*time.Duration(config.IPCacheTimeout))
	if err != nil {
		logger.Error(err)
	}

	configurator := internal.NewDNSConfigurator(config, cache, logger)
	configurator.Configure()
}

func parseCmd() *cmdConfig {
	cmdConfig := &cmdConfig{}
	flag.StringVar(&cmdConfig.ConfigFile, "config", defaultConfigFile, configUsage)
	flag.StringVar(&cmdConfig.ConfigFile, "c", defaultConfigFile, configUsage+" (shorthand)")

	flag.BoolVar(&cmdConfig.Verbose, "verbose", false, verboseUsage)
	flag.BoolVar(&cmdConfig.Verbose, "v", false, verboseUsage+" (shorthand)")

	flag.Parse()

	return cmdConfig
}
