# Overview

Classroom seating is a web application empowering teachers to generate classroom seating plans.

## Meta
- [Homepage](https://github.com/Go-Team-Gamma/classroom-seating)
- [Wiki](https://github.com/Go-Team-Gamma/classroom-seating/wiki)
- [Trello Board](https://trello.com/b/pqQOUld5/classroom-seating)
- [Team Chat](https://nextcloud.groovestomp.com/call/u9dksneb)

## Dependencies
- [Go](https://golang.org/)
- MySQL / MariaDB

## Setup
```
go get github.com/Go-Team-Gamma/classroom-seating
cd $(go env GOPATH)/src/github.com/Go-Team-Gamma/classroom-seating
cp cfg.toml.example cfg.toml
```
Now edit `cfg.toml` appropriately.  
Finally:
```
go build
```

## Running
```
go run *.go # Without building an explicit executable.
./classroom-seating # If using `go build' to generate an executable.
```

## Configuration
- N/A (yet)