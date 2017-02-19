# [nebulo](https://github.com/krostar/nebulo) [![build status](https://travis-ci.org/krostar/nebulo.svg?branch=dev)](https://travis-ci.org/krostar/nebulo)
## Project

## Configuration
The configuration for nebulo's binary can be made in different ways.
First, the configuration manager try to load the configuration from the `/etc/nebulo/config.ini` file, then feom the `./config.ini` and then the command line.

Every new loaded configuration override the previous one, only newly defined properties are overloaded.

### Command line
```sh
$> nebulo -a 10.0.0.1 --port 8080
```

### Files
```INI
environment=dev
address=127.0.0.1
port=17241
logging-file=/var/log/nebulo/log.txt
verbose=debug
```

### Options
```sh
$> nebulo --help
Usage:
  nebulo

Application Options:
  -h, --help                                                show this help message
      --config-gen=                                         generate a configuration file for the actual configuration to the specified file and quit
  -c, --config-file=                                        specify a configuration file (be cautious on infinite-recursive-configuration)
  -e, --environment=[dev|alpha|prod]                        environment to use for external services connection purpose - this parameter is required
  -a, --address=                                            override environment address to use to listen to (default: depend on -e (environment))
  -p, --port=                                               override environment port to use to listen to (default: depend on -e (environment))
      --tls-crt-file=                                       tls certificate file used to encrypt communication
      --tls-key-file=                                       tls certificate key used to encrypt communication
  -l, --logging-file=                                       the file where write the log (default: no file, standart output)
  -v, --verbose=[quiet|critical|error|warning|info|request|debug] level of information to write on standart output or in a file (default: debug)
```


## Documentation
The API documentation and the Golang documentation of this project for the **dev** environment are available here: [doc.nebulo.io/dev](https://doc.nebulo.io/dev)

## Contribute to the project
### Before you started
#### Check your golang installation
Make sure `golang` is installed and is at least in version **1.7** and your `$GOPATH` environment variable set in your working directory
```sh
$> go version
go version go1.7.4 linux/amd64
$> echo $GOPATH
/home/krostar/go
```

If you dont have `golang` installed or if your `$GOPATH` environment variable isn't set, please visit [Golang: Getting Started](https://golang.org/doc/install) and [Golang: GOPATH](https://golang.org/doc/code.html#GOPATH)

> It may be a good idea to add `$GOPATH/bin` and `$GOROOT/bin` in your `$PATH` environment !

#### Download the project
```sh
# Traditional way
$> mkdir -p $GOPATH/src/github.com/krostar/
$> git -c $GOPATH/src/github.com/krostar/ clone https://github.com/krostar/nebulo.git

# Efficient way
$> go get github.com/krostar/nebulo
```

#### Download the dependencies manager
```sh
$> go get -u github.com/kardianos/govendor
```

#### Use our Makefile
We are using a Makefile to everything we need (build, release, tests, documentation, ...).
```sh
# Build the project (by default generated binary will be in <root>/build/bin/nebulo)
$> make
# Run the project
$> make run
# Run the project with command line option
$> make ARGS="--environment dev" run
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
- comment on each 7lines+ functions
- inline comments beginning with a lowercase caracter

Make sure to use `make test` before submitting a merge request !

### Other things
- don't commit dependencies (see [.vendor/vendor.json](https://github.com/kardianos/govendor) configuration file)
- make unit tests for each features!

> In the [atom editor](https://atom.io/) the package `go-plus` is really
> You probably want to use our [`.editorconfig`](http://editorconfig.org) file.

## Report bugs
Create a [issue](https://github.com/krostar/nebulo/issues) or contact [bug[at]nebulo[dot]io](mailto:bug@nebulo.io)
