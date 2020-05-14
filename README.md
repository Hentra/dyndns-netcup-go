# dyndns-netcup-go 
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
 1. Download the lastest [binary](#) for your OS
 2. `cd` to the file you downloaded and unzip
 3. Run `dyndns-netcup-go` as described in [Usage](#usage)

### From source 
First, install [Go](https://golang.org/doc/install) as
recommended.  After that run following commands:

    git clone https://github.com/Hentra/dyndns-netcup-go.git cd 
    dyndns-netcup-go
    go build

This will create a binary named `dyndns-netcup-go` in your current directory.
Refer to [Usage](#usage) for further information.

## Usage
 1. Move/rename the file `example.yml` to `config.yml` and fill out all the
fields. There are some comments in the file for further information. The
filename **has** to be `config.yml`
 2. Run `dyndns-netcup-go` and it will configure your DNS Records.

It might be necessary to run this program every few minutes. That interval
depends on how you configured your TTL.

## Contributing 
For any feature requests and or bugs open up an
[Issue](https://github.com/Hentra/dyndns-netcup-go/issues).  Feel free to also
add a pull request and I will have a look on it.

