# OAuth2 Proxy

[![Build Status](http://github.dronehippie.de/api/badges/webhippie/oauth2-proxy/status.svg)](http://github.dronehippie.de/webhippie/oauth2-proxy)
[![Go Doc](https://godoc.org/github.com/webhippie/oauth2-proxy?status.svg)](http://godoc.org/github.com/webhippie/oauth2-proxy)
[![Go Report](http://goreportcard.com/badge/github.com/webhippie/oauth2-proxy)](http://goreportcard.com/report/github.com/webhippie/oauth2-proxy)
[![](https://images.microbadger.com/badges/image/tboerger/oauth2-proxy.svg)](http://microbadger.com/images/tboerger/oauth2-proxy "Get your own image badge on microbadger.com")
[![Join the chat at https://gitter.im/webhippie/general](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/webhippie/general)
[![Stories in Ready](https://badge.waffle.io/webhippie/oauth2-proxy.svg?label=ready&title=Ready)](http://waffle.io/webhippie/oauth2-proxy)

**This project is under heavy development, it's not in a working state yet!**

A reverse proxy and static file server that provides an authentication layer via OAuth2 to any web application that doesn't support it natively.


## Install

You can download prebuilt binaries from the GitHub releases or from our [download site](http://dl.webhippie.de/misc/oauth2-proxy). You are a Mac user? Just take a look at our [homebrew formula](https://github.com/webhippie/homebrew-webhippie). If you are missing an architecture just write us on our nice [Gitter](https://gitter.im/webhippie/general) chat. If you find a security issue please contact thomas@webhippie.de first.


## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). As this project relies on vendoring of the dependencies and we are not exporting `GO15VENDOREXPERIMENT=1` within our makefile you have to use a Go version `>= 1.6`. It is also possible to just simply execute the `go get github.com/webhippie/oauth2-proxy` command, but we prefer to use our `Makefile`:

```bash
go get -d github.com/webhippie/oauth2-proxy
cd $GOPATH/src/github.com/webhippie/oauth2-proxy
make clean build

./oauth2-proxy -h
```


## Contributing

Fork -> Patch -> Push -> Pull Request


## Authors

* [Thomas Boerger](https://github.com/tboerger)


## License

Apache-2.0


## Copyright

```
Copyright (c) 2017 Thomas Boerger <thomas@webhippie.de>
```
