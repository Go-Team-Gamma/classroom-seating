# Overview

Classroom seating is a web application empowering teachers to generate classroom seating plans.

## Meta
- [Homepage](https://github.com/Go-Team-Gamma/classroom-seating)
- [Wiki](https://github.com/Go-Team-Gamma/classroom-seating/wiki)
- [Trello Board](https://trello.com/b/pqQOUld5/classroom-seating)
- [Team Chat](https://nextcloud.groovestomp.com/call/u9dksneb)

## Dependencies
- PostgreSQL
- [Go](https://golang.org/) 1.11
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cli)

## Setup
```
go get -u github.com/Go-Team-Gamma/classroom-seating
cd $(go env GOPATH)/src/github.com/Go-Team-Gamma/classroom-seating
cp cfg.toml.example cfg.toml
cp .env.example .env
```
Now edit `cfg.toml` and `.env` appropriately.
`.env` isn't strictly required. It provides project-local environment variables.

Finally:
```
. .env

PGUSER=postgres psql template1 -c "CREATE USER $PGUSER WITH PASSWORD '<PASSWORD>';"
PGUSER=postgres psql template1 -c "CREATE DATABASE $PGDATABASE;"
PGUSER=postgres psql testdb -c 'CREATE EXTENSION "uuid-ossp";'
PGUSER=postgres psql testdb -c "CREATE EXTENSION pgcrypto;"
PGUSER=postgres psql testdb -c "CREATE EXTENSION chkpass;"
PGUSER=postgres psql testdb -c "GRANT ALL PRIVILEGES ON DATABASE testdb TO $PGUSER;"

bin/migrate up

go build
```

## Running
```
go run *.go # Without building an explicit executable.
./classroom-seating # If using `go build' to generate an executable.
```

## Configuration
- N/A (yet)