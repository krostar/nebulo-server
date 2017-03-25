# [nebulo](https://github.com/krostar/nebulo) [![License](https://img.shields.io/github/license/krostar/nebulo.svg)](https://tldrlegal.com/license/gnu-general-public-license-v3-(gpl-3)) [![GitHub release](https://img.shields.io/github/release/krostar/nebulo.svg)](https://github.com/krostar/nebulo/releases/latest) [![Godoc](https://godoc.org/github.com/krostar/nebulo?status.svg)](https://godoc.org/github.com/krostar/nebulo)

Nebulo is a secure way of instant messaging that respect and protect your privacy.

[![Build status](https://travis-ci.org/krostar/nebulo.svg?branch=dev)](https://travis-ci.org/krostar/nebulo) [![Go Report Card](https://goreportcard.com/badge/github.com/krostar/nebulo)](https://goreportcard.com/report/github.com/krostar/nebulo) [![Codebeat status](https://codebeat.co/badges/0d3bbf0b-9c5b-44b2-95ae-d29438c89730)](https://codebeat.co/projects/github-com-krostar-nebulo-dev) [![Coverage status](https://coveralls.io/repos/github/krostar/nebulo/badge.svg?branch=dev)](https://coveralls.io/github/krostar/nebulo?branch=dev)

## Usage
The configuration for nebulo's binary can be made in different ways.
First, the configuration manager try to load the configuration from the `/etc/nebulo/config.ini` file, then from the `./config.ini` and then the command line.

Every new loaded configuration override the previous one, only newly defined properties are overloaded.

### Options
```sh
$> nebulo --help
Usage:
  nebulo

Application Options:
  -h, --help                                                      show this help message
      --config-gen=                                               generate a configuration file for the actual configuration to the specified file and quit

Configuration Options:
      --config-dont-load-default                                  choose to load or not the default configuration files (default: false)
  -c, --config-file=                                              specify a configuration file (be cautious on infinite-recursive-configuration)

Environment Options:
  -e, --environment=[dev|beta|prod]                               environment to use for external services connection purpose - this parameter is required
  -a, --address=                                                  override environment address to use to listen to (default: depend on -e (environment))
  -p, --port=                                                     override environment port to use to listen to (default: depend on -e (environment))

Logging Options:
  -l, --logging-file=                                             the file where write the log (default: no file, standart output)
  -v, --verbose=[quiet|critical|error|warning|info|request|debug] level of information to write on standart output or in a file (default: debug)

TLS Options:
      --tls-crt-file=                                             tls certificate file used to encrypt communication - this parameter is required
      --tls-key-file=                                             tls certificate key used to encrypt communication - this parameter is required
      --tls-clients-ca-cert-file=                                 tls certification authority used to validate clients certificate for the tls mutual authentication - this parameter is required
      --tls-clients-ca-key-file=                                  tls certification authority key file used to validate clients certificate for the tls mutual authentication - this parameter is required
      --tls-clients-ca-key-pwd=                                   tls certification authority key password used to validate clients certificate for the tls mutual authentication (default: no password)

User Options:
      --user-provider=[file]                                      provider to use to get users informations
      --user-provider-file=                                       provider file path where users informations are stored

```


## Documentation
The API documentation of this project for the **dev** environment is available on [doc.nebulo.io/dev](https://doc.nebulo.io/dev), the Golang documentation is available on the [godoc website](https://godoc.org/github.com/krostar/nebulo)

## Licence
Distributed under GPL-3 License, please see license file, and/or browse [tldrlegal.com](https://tldrlegal.com/license/gnu-general-public-license-v3-(gpl-3)) for more details.

## Contribute to the project
### Report bugs
Create an [issue](https://github.com/krostar/nebulo/issues) or contact [bug[at]nebulo[dot]io](mailto:bug@nebulo.io)

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
$> git -c $GOPATH/src/github.com/krostar/ clone https://github.com/krostar/nebulo.git

# or via go get
$> go get github.com/krostar/nebulo
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
# ...
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
