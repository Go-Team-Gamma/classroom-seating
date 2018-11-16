package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
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

	dbString := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sqlx.Connect("postgres", dbString)
	if err != nil {
		log.Fatalln(err)
	}

	newUser := User{
		Username: "aoman",
		Password: "P4ssword!!",
	}
	stmt, err := db.PrepareNamed(`INSERT INTO  users (username, password) VALUES (:username, :password)`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = stmt.Exec(newUser)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", newUser)

	err = db.Get(&newUser, "SELECT * FROM users WHERE username=$1", newUser.Username)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", newUser)

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
	Id       sql.NullString
	Username string
	Password string
}
