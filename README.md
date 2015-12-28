# dispatch [![Build Status](https://travis-ci.org/khlieng/dispatch.svg?branch=master)](https://travis-ci.org/khlieng/dispatch)

####[Try it!](http://dispatch.khlieng.com/)

### Features
* Searchable history
* Persistent connections
* Multiple users

## Usage
There is a few different ways of getting it:

#### 1. Binary
There will be binary releases.

#### 2. Go
This requires a [Go environment](http://golang.org/doc/install).

```bash
go get github.com/khlieng/dispatch

dispatch
```

To get some help run:
```bash
dispatch help
```

#### 3. Docker
```bash
docker run -p 8080:8080 khlieng/dispatch
```

## Build

### Server
```bash
cd $GOPATH/src/github.com/khlieng/dispatch

go install
```

### Client
This requires [Node.js](https://nodejs.org/download/).

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

For development with hot reloading enabled just run:
```bash
gulp
```
