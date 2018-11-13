package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gocraft/dbr"
	_ "github.com/lib/pq"
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

	dbString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	fmt.Printf("DSN String: %v\n", dbString)
	conn, err := dbr.Open("postgres", dbString, nil)
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
	Port     uint
	Name     string
	User     string
	Password string
}

type HttpServerConfig struct {
	Port uint
}
