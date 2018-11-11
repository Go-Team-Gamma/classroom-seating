package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"html"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var (
		cfg Config
	)

	dat, err := ioutil.ReadFile("cfg.toml")
	if err != nil {
		log.Fatal(err)
	}
	strDat := string(dat)

	if _, err := toml.Decode(strDat, &cfg); err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)

	dbString := fmt.Sprintf("%v:%v@/%v", cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	conn, err := dbr.Open("mysql", dbString, nil)
	if err != nil {
		log.Fatal(err)
	}
	sess := conn.NewSession(nil)
	tx, err := sess.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.RollbackUnlessCommitted()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})

	err = http.ListenAndServe(fmt.Sprintf(":%v", cfg.Server.Port), nil)
	if err != nil {
		log.Fatal("http.listenAndServe: ", err)
	}
}

type Config struct {
	Database DatabaseConfig   `toml:"database"`
	Server   HttpServerConfig `toml:"http_server"`
}

type DatabaseConfig struct {
	Host     string
	Name     string
	User     string
	Password string
}

type HttpServerConfig struct {
	Port uint
}
