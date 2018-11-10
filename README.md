# dispatch [![Build Status](https://travis-ci.org/khlieng/dispatch.svg?branch=master)](https://travis-ci.org/khlieng/dispatch)

#### [Try it!](https://dispatch.khlieng.com)

![Dispatch](https://khlieng.com/dispatch.png?1)

### Features

- Searchable history
- Persistent connections
- Multiple servers and users
- Automatic HTTPS through Let's Encrypt
- Client certificates

## Usage

There is a few different ways of getting it:

### 1. Binary

- **[Windows (x64)](https://github.com/khlieng/dispatch/releases/download/v0.4/dispatch_windows_amd64.zip)**
- **[OS X (x64)](https://github.com/khlieng/dispatch/releases/download/v0.4/dispatch_darwin_amd64.zip)**
- **[Linux (x64)](https://github.com/khlieng/dispatch/releases/download/v0.4/dispatch_linux_amd64.tar.gz)**
- [Other versions](https://github.com/khlieng/dispatch/releases)

### 2. Go

This requires a [Go environment](http://golang.org/doc/install), version 1.10 or greater.

Fetch, compile and run dispatch:

```bash
go get github.com/khlieng/dispatch
dispatch
```

To get some help run:

```bash
dispatch help
```

### 3. Docker

```bash
docker run -p <http port>:80 -p <https port>:443 -v <path>:/data khlieng/dispatch
```

## Build

### Server

```bash
cd $GOPATH/src/github.com/khlieng/dispatch
go install
```

### Client

This requires [Node.js](https://nodejs.org) and [yarn](https://yarnpkg.com).

Fetch the dependencies:

```bash
go get github.com/jteeuwen/go-bindata/...
yarn global add gulp@next
cd $GOPATH/src/github.com/khlieng/dispatch/client
yarn
```

Run the build:

```bash
gulp build
```

The server needs to be rebuilt to embed new client builds.

For development with hot reloading start the frontend:

```bash
gulp
```

And then the backend in a separate terminal:

```bash
dispatch --dev
```

## Libraries

The libraries this project is built with.

### Server

- [Bolt](https://github.com/boltdb/bolt)
- [Bleve](https://github.com/blevesearch/bleve)
- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
- [Lego](https://github.com/xenolf/lego)

### Client

- [React](https://github.com/facebook/react)
- [Redux](https://github.com/reactjs/redux)
- [Immer](https://github.com/mweststrate/immer)
- [react-window](https://github.com/bvaughn/react-window)
- [Lodash](https://github.com/lodash/lodash)

## Big Thanks

Cross-browser Testing Platform and Open Source <3 Provided by [Sauce Labs][homepage]

[homepage]: https://saucelabs.com
