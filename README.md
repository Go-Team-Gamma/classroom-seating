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

```
sudo su postgres # Login as postgres user.
```

Export `$PGNEWUSER`, `$PGPASSWORD` and `$PGDATABASE` according to what you configured in `cfg.toml`:
```
psql template1 -c "CREATE USER $PGNEWUSER WITH PASSWORD '$PGPASSWORD';"
psql template1 -c "CREATE DATABASE $PGDATABASE;"
psql $PGDATABASE -c 'CREATE EXTENSION "uuid-ossp";'
psql $PGDATABASE -c "GRANT ALL PRIVILEGES ON DATABASE $PGDATABASE TO $PGNEWUSER;"

exit # Logout as postgres user.
```

Now do the migrations:
```
bin/db-migrate up
```

And now you can build the project:
```
go build
```

## Running
```
./classroom-seating # If using `go build' to generate an executable.
```

## Configuration
- N/A (yet)

## Testing
```
cp cfg.toml.example cfg.test.toml
```
Now edit `cfg.test.toml`.

```
sudo su postgres # Login as postgres user.
```

Export `$PGNEWUSER`, `$PGPASSWORD` and `$PGDATABASE` according to what you configured in `cfg.test.toml`:
```
psql template1 -c "CREATE USER $PGNEWUSER WITH PASSWORD '$PGPASSWORD';"
psql template1 -c "CREATE DATABASE $PGDATABASE;"
psql $PGDATABASE -c 'CREATE EXTENSION "uuid-ossp";'
psql $PGDATABASE -c "GRANT ALL PRIVILEGES ON DATABASE $PGDATABASE TO $PGNEWUSER;"

exit # Logout as postgres user.
```

Now do the migrations:
```
CFG=cfg.test.toml bin/db-migrate up
```

And finally, you can run tests:
```
go test ./...
```