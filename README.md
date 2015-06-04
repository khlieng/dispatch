# name_pending

####[Try it!](http://np.khlieng.com/)

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
go get github.com/khlieng/name_pending

name_pending
```

To get some help run:
```bash
name_pending help
```

#### 3. Docker
```bash
docker run -p 8080:8080 khlieng/name_pending
```

## Build

### Server
```bash
cd $GOPATH/src/github.com/khlieng/name_pending

go install
```

### Client
This requires [Node.js](https://nodejs.org/download/).

Fetch the dependencies:
```bash
npm install -g gulp
go get github.com/jteeuwen/go-bindata/...
cd $GOPATH/src/github.com/khlieng/name_pending/client
npm install
```

Run the build:
```bash
gulp -p
```

The server needs to be rebuilt after this. For development dropping the -p flag 
will turn off minification and embedding, requiring only one initial server rebuild.
