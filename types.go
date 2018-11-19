package main

import (
	"database/sql"
)

type Config struct {
	Database DatabaseConfig   `toml:"database"`
	Server   HttpServerConfig `toml:"http_server"`
}

type DatabaseConfig struct {
	Host     string
	Port     uint
	Name     string
	User     string
	Password string
}

type HttpServerConfig struct {
	Port uint
}

type User struct {
	Id       sql.NullString
	Username string
	Password string
}
