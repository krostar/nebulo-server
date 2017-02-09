# [nebulo](https://github.com/krostar/nebulo) [![build status](https://travis-ci.org/krostar/nebulo.svg?branch=dev)](https://travis-ci.org/krostar/nebulo)
## Project

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

If you dont have `golang` installed or if your `$GOPATH` environment variable isn't set, please visit [Golang: Getting Started](https://golang.org/doc/install)

> It may be a good idea to add `$GOPATH/bin` and `$GOROOT/bin` in your `$PATH` environment !

#### Download the project
```sh
# Traditional way
$> mkdir -p $GOPATH/src/github.com/krostar/
$> git -c $GOPATH/src/github.com/krostar/ clone https://github.com/krostar/nebulo.git

# Efficient way
$> go get github.com/krostar/nebulo
```

#### Use our Makefile
We are using a Makefile to everything we need (build, release, tests, documentation, ...).
```sh
# Build the project
$> make
# Test the project
$> make test
# Generate documentation
$> make doc
# ...
```


### Guidelines
#### Coding standart
Please, make sure your favorite editor is configured for this project. The source code should be:
- well formatted (`gofmt` (usage of tabulation, no trailing whitespaces, trailing line at the end of the file, ...))
- linter free (`gometalinter --config=.gometalinter.json ./...`)
- comment on each 7lines+ functions
- inline comments beginning with a lowercase caracter

To avoid having a merge request rejected because of this, please use our [`.editorconfig`](http://editorconfig.org) file.

### Other things
- don't commit dependencies (see [.vendor/vendor.json](https://github.com/kardianos/govendor) configuration file)
- make unit tests for each features!

> In the [atom editor](https://atom.io/) the package `go-plus` is really useful

### Report bugs
Create a [issue](https://github.com/krostar/nebulo/issues) or contact [bug[at]nebulo[dot]io](mailto:bug@nebulo.io)
