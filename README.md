# [nebulo-server](https://github.com/krostar/nebulo-server) [![License](https://img.shields.io/github/license/krostar/nebulo.svg)](https://tldrlegal.com/license/gnu-general-public-license-v3-(gpl-3)) [![GitHub release](https://img.shields.io/github/release/krostar/nebulo-server.svg)](https://github.com/krostar/nebulo-server/releases/latest) [![Godoc](https://godoc.org/github.com/krostar/nebulo-server?status.svg)](https://godoc.org/github.com/krostar/nebulo-server)

Nebulo is a secure way of instant messaging that respect and protect your privacy.

[![Build status](https://travis-ci.org/krostar/nebulo-server.svg?branch=dev)](https://travis-ci.org/krostar/nebulo-server) [![Go Report Card](https://goreportcard.com/badge/github.com/krostar/nebulo-server)](https://goreportcard.com/report/github.com/krostar/nebulo-server) [![Codebeat status](https://codebeat.co/badges/65b96953-23e6-4440-8be1-84ccc75e7090)](https://codebeat.co/projects/github-com-krostar-nebulo-server-dev) [![Coverage status](https://coveralls.io/repos/github/krostar/nebulo-server/badge.svg?branch=dev)](https://coveralls.io/github/krostar/nebulo-server?branch=dev)

## Usage
```sh
# check if nebulo is in your $PATH
$>nebulo version

# see commands and parameters
$>nebulo help

# get help on the run command
$>nebulo help run

# copy sample configuration file
$>cp config.sample/json config.json

# fill required values (run `nebulo help run` to know which values are required)
$>vim config.json

# start the server
$>nebulo -c path/to/config.json run
```

## Documentation
The API documentation of this project for the **dev** environment is available on [doc.nebulo.io/dev](https://doc.nebulo.io/dev), the Golang documentation is available on the [godoc website](https://godoc.org/github.com/krostar/nebulo-server)

## Licence
Distributed under GPL-3 License, please see license file, and/or browse [tldrlegal.com](https://tldrlegal.com/license/gnu-general-public-license-v3-(gpl-3)) for more details.

## Contribute to the project
### Report bugs
Create an [issue](https://github.com/krostar/nebulo-server/issues) or contact [bug[at]nebulo[dot]io](mailto:bug@nebulo.io)

### Before you started
#### Check your golang installation
Make sure `golang` is installed and is at least in version **1.8** and your `$GOPATH` environment variable set in your working directory
```sh
$> go version
go version go1.8 linux/amd64
$> echo $GOPATH
/home/krostar/go
```

If you don't have `golang` installed or if your `$GOPATH` environment variable isn't set, please visit [Golang: Getting Started](https://golang.org/doc/install) and [Golang: GOPATH](https://golang.org/doc/code.html#GOPATH)

> It may be a good idea to add `$GOPATH/bin` and `$GOROOT/bin` in your `$PATH` environment!

#### Download the project
```sh
# Manually
$> mkdir -p $GOPATH/src/github.com/krostar/
$> git -c $GOPATH/src/github.com/krostar/ clone https://github.com/krostar/nebulo-server.git

# or via go get
$> go get github.com/krostar/nebulo-server
```

#### Download the tool manager
```sh
$> go get -u github.com/twitchtv/retool
```

#### Use our Makefile
We are using a Makefile to everything we need (build, release, tests, documentation, ...).
```sh
# Get the dependencies and tools
$> make vendor

# Build the project (by default generated binary will be in <root>/build/bin/nebulo)
$> make build

# Run the project without arguments
$> make run

# Run the project with arguments
$> make run ARGS="--environment dev"

# Test the project
$> make test

# Generate documentation
$> make doc

# Generate release
$> make release TAG=1.2.3
```

### Guidelines
#### Coding standart
Please, make sure your favorite editor is configured for this project. The source code should be:
- well formatted (`gofmt` (usage of tabulation, no trailing whitespaces, trailing line at the end of the file, ...))
- linter free (`gometalinter --config=.gometalinter.json ./...`)
- with inline comments beginning with a lowercase caracter

Make sure to use `make test` before submitting a pull request!

### Other things
- use the dependencies manager and update them (see [govendor](https://github.com/kardianos/govendor) and [retool](https://github.com/twitchtv/retool))
- write unit tests
