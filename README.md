# DYNDNS NETCUP GO
![Build](https://github.com/Hentra/dyndns-netcup-go/workflows/Build/badge.svg?branch=master)
[![Issues](https://img.shields.io/github/issues/Hentra/dyndns-netcup-go)](https://github.com/Hentra/dyndns-netcup-go/issues)
[![Release](https://img.shields.io/github/release/Hentra/dyndns-netcup-go?include_prereleases)](https://github.com/Hentra/dyndns-netcup-go/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Hentra/dyndns-netcup-go)](https://goreportcard.com/report/github.com/Hentra/dyndns-netcup-go)

Dyndns client for the netcup DNS API written in go. Not
related to netcup GmbH. It is **heavily** inspired by 
[this](https://github.com/stecklars/dynamic-dns-netcup-api) 
project which might be also a good solution for your 
dynamic DNS needs. 

## Table of Contents
<!-- vim-markdown-toc GFM -->

* [Features](#features)
* [Installation](#installation)
	* [Docker compose](#docker-compose-example-using-secret-files)
	* [Docker CLI](#docker-cli)
  * [Environment Variables](#environment-variables)
	* [Manual](#manual)
	* [From source](#from-source)
* [Usage](#usage)
	* [Prequisites](#prequisites)
	* [Run dyndns-netcup-go](#run-dyndns-netcup-go)
		* [Commandline flags](#commandline-flags)
	* [Cache](#cache)
* [Contributing](#contributing)

<!-- vim-markdown-toc -->

## Features

* Multi domain support
* Subdomain support
* TTL update support
* Creation of a DNS record if it doesn't already exist
* Multi host support (nice when you need to update both `@` and `*`) 
* IPv6 support
* secure Docker support 
* secret files

If you need additional features please open up an
[Issue](https://github.com/Hentra/dyndns-netcup-go/issues).

## Installation 

### Docker compose example using secret files

You need to create config.yml file in the same directory as the docker-compose.yml file, take a look at [config/example.yml](config/example.yml) for an example.
For a Docker setup do not save your secrets (api key, etc.) directly in the config.yml, but rather as environment variable or secret file, as shown below!
For secrets management create three files under `secrets/` with the names `customernr`, `apikey` and `apipassword`, like in the [secrets/](secrets/) directory.
To further protect the secrets from unauthorized access, make sure it is owned by the user that runs dyndns-netcup-go, by default being the UID 62534 and make it read-only:
```shell
sudo chown 62534:62534 secrets/*
sudo chmod 440 secrets/*
```

After that setup you can use the following docker-compose.yml as an example, also available in [docker-compose.yml](docker-compose.yml):
```compose.yml
services:
  dyn-dns:
    image: ghcr.io/hentra/dyndns-netcup-go
    container_name: Netcup-Dyn-DNS
    environment:
      - INTERVAL=300
      - CUSTOMERNR_FILE=/run/secrets/customernr
      - APIKEY_FILE=/run/secrets/apikey
      - APIPASSWORD_FILE=/run/secrets/apipassword
    secrets:
      - customernr
      - apikey
      - apipassword
    volumes:
      - ./config.yml:/config.yml
    security_opt:
      - no-new-privileges
    cap_drop:
      - ALL
    restart: unless-stopped
    networks:
      - ipv6-enabled-network

secrets:
  customernr:
    file: ./secrets/customernr
  apikey:
    file: ./secrets/apikey
  apipassword:
    file: ./secrets/apipassword

networks:
  ipv6-enabled-network:
    enable_ipv6: true
```

### Docker CLI

    docker run -d \
        -v $(pwd)/config.yml:/config.yml \
        -e INTERVAL=300 \
        -e CUSTOMERNR=111111 \
        -e APIKEY=my-fancy-api-key \
        -e APIPASSWORD=my-fancy-api-pw \
        ghcr.io/hentra/dyndns-netcup-go

This allows you to store the configuration in plain text(e.g. git) and inject the secrets safely from a secret management solution.

### Environment Variables

| Environment Variable | Description                                                             |
|----------------------|-------------------------------------------------------------------------|
| INTERVAL             | defines the interval of DNS updates in seconds                          |
| CUSTOMERNR_FILE      | location of the secrets file containing your customer number from netcup|
| APIKEY_FILE          | location of the secrets file containing your API key from netcup        |
| APIPASSWORD_FILE     | location of the secrets file containing your API password from netcu    |
| CUSTOMERNR           | alternative to secrets file: customer number from netcup                |
| APIKEY               | alternative to secrets file: containing your API key from netcup        |
| APIPASSWORD          | alternative to secrets file: containing your API password from netcu    |

### Manual
 1. Download the lastest [binary](https://github.com/Hentra/dyndns-netcup-go/releases) for your OS
 2. `cd` to the file you downloaded and unzip
 3. Put `dyndns-netcup-go` somewhere in your path

### From source 
First, install [Go](https://golang.org/doc/install) as
recommended.  After that run following commands:

    git clone https://github.com/Hentra/dyndns-netcup-go.git 
    cd dyndns-netcup-go
    go install

This will create a binary named `dyndns-netcup-go` and install it to your go
binary home. Make sure your `GOPATH` environment variable is set. 

Refer to [Usage](#usage) for further information.

## Usage

### Prequisites

* You need to have a netcup account and a domain, obviously.
* Then you need an apikey and apipassword.
  [Here](https://www.netcup-wiki.de/wiki/CCP_API#Authentifizierung) is a
description (in German) on how you get those.

### Run dyndns-netcup-go
1. Move/rename the [example configuration](./config/example.yml) `config/example.yml` 
to `config.yml` and fill out all the fields. There are some comments in the file for further information. 
2. Run `dyndns-netcup-go -v` in the **same** directory as your configuration file and it will
configure your DNS Records. You can specify the location of the
configuration file with the `-c` or `-config` flag if you don't want to run
it in the same directory. To disable the output for information remove the `-v` flag. You will
still get the output from errors.

It might be necessary to run this program every few minutes. That interval
depends on how you configured your TTL.

#### Commandline flags
For a list of all available command line flags run `dyndns-netcup-go -h`.

### Cache
Without the cache the application would lookup its ip addresses and fetch the DNS
records from netcup. After that it will compare the specified hosts in the DNS
records with the current ip addresses and update if necessary. 

As reported in [this issue](https://github.com/Hentra/dyndns-netcup-go/issues/1)
it would be also possible to store the ip addresses between two runs of the
application and only fetch DNS records from netcup when they differ. 

To enable the cache configure the two variables `IP-CACHE` and
`IP-CACHE-TIMEOUT` as according to the comments in `example.yml`.

## Contributing 
For any feature requests and or bugs open up an
[Issue](https://github.com/Hentra/dyndns-netcup-go/issues).  Feel free to also
add a pull request and I will have a look on it.

