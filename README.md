# unused [![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godoc] [![Travis](https://img.shields.io/travis/gostaticanalysis/unused.svg?style=flat-square)][travis] [![Go Report Card](https://goreportcard.com/badge/github.com/gostaticanalysis/unused)](https://goreportcard.com/report/github.com/gostaticanalysis/unused) [![codecov](https://codecov.io/gh/gostaticanalysis/unused/branch/master/graph/badge.svg)](https://codecov.io/gh/gostaticanalysis/unused)

`unused` finds unused package level identifiers.

## Install

```sh
$ go get github.com/gostaticanalysis/unused
```

## Usage

```sh
$ go vet -vettool=`which unused` pkgname
```

<!-- links -->
[godoc]: http://godoc.org/github.com/gostaticanalysis/unused
[travis]: https://travis-ci.org/gostaticanalysis/unused
