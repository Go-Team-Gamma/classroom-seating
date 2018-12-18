package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/husobee/vestigo"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	db *sqlx.DB
)

func main() {
	router := vestigo.NewRouter()

	dat, err := ioutil.ReadFile("cfg.toml")
	if err != nil {
		log.Fatal(err)
	}
	strDat := string(dat)

	var cfg Config
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

	db, err = sqlx.Connect("postgres", dbString)
	if err != nil {
		log.Fatalln(err)
	}

	router.Get("/", ShowRoot)
	router.Get("/register", ShowRegistration)
	router.Get("/login", ShowLogin)
	router.Post("/users", CreateUser)
	router.Post("/login", Login)
	router.Get("/logout", Logout)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Server.Port), router))
}
