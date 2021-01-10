package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	CustomerNumber int      `yaml:"CUSTOMERNR"`
	ApiKey         string   `yaml:"APIKEY"`
	ApiPassword    string   `yaml:"APIPASSWORD"`
	IPCache        string   `yaml:"IP-CACHE"`
	IPCacheTimeout int      `yaml:"IP-CACHE-TIMEOUT"`
	Domains        []Domain `yaml:"DOMAINS"`
}

type Domain struct {
	Name  string   `yaml:"NAME"`
	IPv6  bool     `yaml:"IPV6"`
	IPv4  bool     `yaml:"IPV4"`
	TTL   int      `yaml:"TTL"`
	Hosts []string `yaml:"HOSTS"`
}

func LoadConfig(filename string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (d *Domain) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawDomain Domain
	raw := rawDomain{
		IPv4: true,
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*d = Domain(raw)
	return nil
}
