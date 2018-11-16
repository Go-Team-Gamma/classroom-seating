package main

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gocraft/dbr"
	_ "github.com/lib/pq"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
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

	dbString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	conn, err := dbr.Open("postgres", dbString, nil)
	if err != nil {
		log.Fatal(err)
	}
	sess := conn.NewSession(nil)

	newUser := &User{
		Username: "aoman",
		Password: "P4ssword!!",
	}
	_, err = sess.InsertInto("users").
		Columns("username", "password").
		Record(newUser).
		Exec()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", newUser.Id)

	var users []User
	sess.Select("*").From("users").Load(&users)

	fmt.Printf("users: %+v\n", users)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}

		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}

		viewModel := struct {
			UrlPath string
		}{
			UrlPath: html.EscapeString(r.URL.Path),
		}

		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			log.Fatal(err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, viewModel)
		fmt.Fprintf(w, buf.String())
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		fmt.Printf("%+v\n", r)
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

type User struct {
	Id       dbr.NullString
	Username string
	Password string
}
