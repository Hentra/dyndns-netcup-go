package internal

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents a config.
type Config struct {
	CustomerNumber int      `yaml:"CUSTOMERNR"`
	APIKey         string   `yaml:"APIKEY"`
	APIPassword    string   `yaml:"APIPASSWORD"`
	IPCache        string   `yaml:"IP-CACHE"`
	IPCacheTimeout int      `yaml:"IP-CACHE-TIMEOUT"`
	IPv4Service    string   `yaml:"IPV4-SERVICE"`
	IPv6Service    string   `yaml:"IPV6-SERVICE"`
	Domains        []Domain `yaml:"DOMAINS"`
}

// Domain represents a domain.
type Domain struct {
	Name  string   `yaml:"NAME"`
	IPv6  bool     `yaml:"IPV6"`
	IPv4  bool     `yaml:"IPV4"`
	TTL   int      `yaml:"TTL"`
	Hosts []string `yaml:"HOSTS"`
}

// LoadConfig returns a config loaded from a specified location. It will
// return an error if there is no file in the specified location or it is
// unable to read it.
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

// UnmarshalYAML is implemented to override the default value of
// the IPv4 field of a Domain with true.
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

// CacheEnabled returns whether the cache is enabled in the
// configuration.
func (c *Config) CacheEnabled() bool {
	return c.IPCacheTimeout > 0
}

// IPv6Enabled returns true if at least one domain needs the AAAA
// record configured.
func (c *Config) IPv6Enabled() bool {
	for _, domain := range c.Domains {
		if domain.IPv6 {
			return true
		}
	}

	return false
}

// IPv4Enabled returns true if at least one domain needs the A
// record configured.
func (c *Config) IPv4Enabled() bool {
	for _, domain := range c.Domains {
		if domain.IPv4 {
			return true
		}
	}

	return false
}
