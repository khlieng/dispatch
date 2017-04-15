# dispatch [![Build Status](https://travis-ci.org/khlieng/dispatch.svg?branch=master)](https://travis-ci.org/khlieng/dispatch)

#### [Try it!](https://dispatch.khlieng.com)

![Dispatch](https://khlieng.com/dispatch.png)

### Features
* Searchable history
* Persistent connections
* Multiple servers and users
* Automatic HTTPS through Let's Encrypt
* Client certificates

## Usage
There is a few different ways of getting it:

### 1. Binary
- **[Windows (x64)](https://github.com/khlieng/dispatch/releases/download/v0.2/dispatch_windows_amd64.zip)**
- **[OS X (x64)](https://github.com/khlieng/dispatch/releases/download/v0.2/dispatch_darwin_amd64.zip)**
- **[Linux (x64)](https://github.com/khlieng/dispatch/releases/download/v0.2/dispatch_linux_amd64.tar.gz)**
- [Other versions](https://github.com/khlieng/dispatch/releases)

### 2. Go
This requires a [Go environment](http://golang.org/doc/install), version 1.8 or greater.

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
This requires [Node.js](https://nodejs.org).

Fetch the dependencies:
```bash
npm install -g gulp
go get github.com/jteeuwen/go-bindata/...
cd $GOPATH/src/github.com/khlieng/dispatch/client
npm install
```

Run the build:
```bash
gulp build
```

The server needs to be rebuilt after this.

For development with hot reloading enabled run:
```bash
gulp
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
- [React Router](https://github.com/ReactTraining/react-router)
- [React Virtualized](https://github.com/bvaughn/react-virtualized)
- [Immutable](https://github.com/facebook/immutable-js)
- [Lodash](https://github.com/lodash/lodash)
