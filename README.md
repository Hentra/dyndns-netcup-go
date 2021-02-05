# DYNDNS NETCUP GO
![Build](https://github.com/Hentra/dyndns-netcup-go/workflows/Build/badge.svg?branch=master)
[![Issues](https://img.shields.io/github/issues/Hentra/dyndns-netcup-go)](https://github.com/Hentra/dyndns-netcup-go/issues)
[![Release](https://img.shields.io/github/release/Hentra/dyndns-netcup-go?include_prereleases)](https://github.com/Hentra/dyndns-netcup-go/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Hentra/dyndns-netcup-go)](https://goreportcard.com/report/github.com/Hentra/dyndns-netcup-go)

Dyndns client for the netcup dns API written in go. Not
related to netcup GmbH. It is **heavily** inspired by 
[this](https://github.com/stecklars/dynamic-dns-netcup-api) 
project which might be also a good solution for your 
dynamic dns needs. 

## Table of Contents
<!-- vim-markdown-toc GFM -->

* [Features](#features)
  * [Implemented](#implemented)
  * [Missing](#missing)
* [Installation](#installation)
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

### Implemented
* Multi domain support
* Subdomain support
* TTL update support
* Creation of a DNS record if it doesn't already exist.
* Multi host support (nice when you need to update both `@` and `*`) 
* IPv6 support
* Verbose option (when you specify `-v` you get plenty information)
* Cache IP of the executing machine and only update when it changes

### Missing

* MX entry support

There are currently no plans to implement this features. If you need those (or additional features) please
open up an [Issue](https://github.com/Hentra/dyndns-netcup-go/issues).

## Installation 

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
1. Move/rename the file `example.yml` to `config.yml` and fill out all the
fields. There are some comments in the file for further information. 
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
`IP-CACHE-LOCATION` as according to the comments in `example.yml`.

## Contributing 
For any feature requests and or bugs open up an
[Issue](https://github.com/Hentra/dyndns-netcup-go/issues).  Feel free to also
add a pull request and I will have a look on it.

