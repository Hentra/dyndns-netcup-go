# DYNDNS NETCUP GO
![Build](https://github.com/Hentra/dyndns-netcup-go/workflows/Build/badge.svg?branch=master)
![Issues](https://img.shields.io/github/issues/Hentra/dyndns-netcup-go)
![Release](https://img.shields.io/github/release/Hentra/dyndns-netcup-go?include_prereleases)

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
* [Contributing](#contributing)

<!-- vim-markdown-toc -->

## Features

### Implemented
* Multi domain support
* Subdomain support
* TTL update support
* Creation of a DNS record if it doesn't already exists.
* Multi host support (nice when you need to update both `@` and `*`) 

### Missing
* IPv6 support
* Quiet option (output is always really verbose)

## Installation 

### Manual
 1. Download the lastest [binary](https://github.com/Hentra/dyndns-netcup-go/releases) for your OS
 2. `cd` to the file you downloaded and unzip
 3. Put `dyndns-netcup-go` somewhere in your path

### From source 
First, install [Go](https://golang.org/doc/install) as
recommended.  After that run following commands:

    git clone https://github.com/Hentra/dyndns-netcup-go.git cd 
    dyndns-netcup-go
    go install

This will create a binary named `dyndns-netcup-go` and install it to your go binary home.
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
2. Run `dyndns-netcup-go` in the **same** directory as your configuration file and it will
configure your DNS Records. You can specify the location of the
configuration file with the `-c` or `-config` flag if you dont want to run
it in the same directory.

It might be necessary to run this program every few minutes. That interval
depends on how you configured your TTL.

## Contributing 
For any feature requests and or bugs open up an
[Issue](https://github.com/Hentra/dyndns-netcup-go/issues).  Feel free to also
add a pull request and I will have a look on it.

