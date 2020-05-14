package main

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type Config struct {
    CustomerNumber int `yaml:"CUSTOMERNR"`
    ApiKey string `yaml:"APIKEY"`
    ApiPassword string `yaml:"APIPASSWORD"`
    Domains []Domain `yaml:"DOMAINS"`
}

type Domain struct {
    Name string `yaml:"NAME"`
    IPv6 bool `yaml:"IPV6"`
    TTL int `yaml:"TTL"`
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
