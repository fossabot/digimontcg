# digimontcg
They are the champions

## Setup
This project depends on the [Go programming language](https://golang.org/dl/).
I like to use a [POSIX-compatible Makefile](https://pubs.opengroup.org/onlinepubs/9699919799.2018edition/utilities/make.html) to facilitate the various project operations but traditional [Go commands](https://pkg.go.dev/cmd/go) will work just as well.

## Building
To build the application into a standalone binary, run:
```
make
```

## Local Development
### Running
To start the web server:
```
make web
```

### Testing
Unit and integration tests can be ran after starting the aforementioned services:
```
make test
```
