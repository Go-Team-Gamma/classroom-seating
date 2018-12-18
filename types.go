package main

import (
	"database/sql"
	"time"
	"net/http"
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
	Id        sql.NullString
	Username  string
	Password  []byte
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

type PageData struct {
	Authenticated bool
	Title         string
}

type Cookier interface {
	Cookie(name string) (*http.Cookie, error)
}
