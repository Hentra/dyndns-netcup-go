package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Hentra/dyndns-netcup-go/internal"
)

const (
	configFileLocation = "/config.yml"
	ipCacheLocation    = "/ipcache"
	defaultInterval    = time.Minute
	intervalEnv        = "INTERVAL"
)

func main() {
	interval, err := parseEnv()
	if err != nil {
		log.Fatal("Could not parse interval: ", err)
	}

	logger := internal.NewLogger(true)

	config, err := internal.LoadConfig(configFileLocation)
	if err != nil {
		logger.Error("Error loading config file. Make sure a config file is mounted to ", configFileLocation, ":", err)
	}

	config.IPCache = ipCacheLocation

	cache, err := internal.NewCache(config.IPCache, time.Second*time.Duration(config.IPCacheTimeout))
	if err != nil {
		logger.Error(err)
	}

	configurator := internal.NewDNSConfigurator(config, cache, logger)
	for {
		logger.Info("configure DNS records")
		configurator.Configure()
		time.Sleep(interval)
	}
}

func parseEnv() (time.Duration, error) {
	intervalTime := defaultInterval
	if interval, exists := os.LookupEnv(intervalEnv); exists {
		intervalSeconds, err := strconv.Atoi(interval)
		if err != nil {
			return 0, err
		}
		intervalTime = time.Duration(intervalSeconds) * time.Second
	}

	return intervalTime, nil
}
