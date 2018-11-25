package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func cleanupDB() {
	tx := db.MustBegin()
	tx.MustExec("DELETE FROM authentications")
	tx.MustExec("DELETE FROM users")
	tx.Commit()
}

func TestMain(m *testing.M) {
	dat, err := ioutil.ReadFile("cfg.test.toml")
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	strDat := string(dat)

	var cfg Config
	if _, err := toml.Decode(strDat, &cfg); err != nil {
		log.Fatalln(err)
		os.Exit(1)
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
		os.Exit(1)
	}

	cleanupDB()
	os.Exit(m.Run())
}

func TestLogHandlerIntro(t *testing.T) {
	data := url.Values{}
	data.Set("name", "thing")
	data.Add("colors", "blue")
	data.Add("colors", "red")
	logHandlerIntro("POST", "/", data)
	// Output: POST "/": map[]
}

type GoodCookier struct {
	Value string
}

func (c GoodCookier) Cookie(string) (*http.Cookie, error) {
	return &http.Cookie{Name: "test", Value: c.Value}, nil
}

type BadCookier struct{}

func (BadCookier) Cookie(string) (*http.Cookie, error) {
	return nil, errors.New("Testing Error")
}

func TestAuthenticate(t *testing.T) {
	var userId string
	var authToken string
	tx := db.MustBegin()
	tx.MustExec(`INSERT INTO users (username, password) VALUES ($1, $2)`, `testuser`, `testpassword`)
	tx.Get(&userId, `SELECT id FROM users WHERE username = $1`, `testuser`)
	tx.MustExec(`INSERT INTO authentications (user_id) VALUES ($1)`, userId)
	tx.Get(&authToken, `SELECT token FROM authentications WHERE user_id = $1`, userId)
	tx.Commit()
	defer cleanupDB()

	goodCookie := GoodCookier{Value: authToken}
	auth, err := authenticate(goodCookie)
	if err != nil {
		t.Errorf("Expected err to be nil, but was '%v'\n", err)
	}
	if auth == "" {
		t.Error("Expected auth string not to be empty")
	}
}

func TestAuthenticate_missingCookie(t *testing.T) {
	_, err := authenticate(BadCookier{})
	if err == nil {
		t.Fail()
	}
	if err.Error() != "Testing Error" {
		t.Errorf("%v != \"Testing Error\"\n", err)
	}
}

func TestAuthenticate_missingAuthentication(t *testing.T) {
	var userId string
	var authToken string
	tx := db.MustBegin()
	tx.MustExec(`INSERT INTO users (username, password) VALUES ($1, $2)`, `testuser`, `testpassword`)
	tx.Get(&userId, `SELECT id FROM users WHERE username = $1`, `testuser`)
	tx.MustExec(`INSERT INTO authentications (user_id) VALUES ($1)`, userId)
	tx.Get(&authToken, `SELECT token FROM authentications WHERE user_id = $1`, userId)
	tx.MustExec(`UPDATE authentications SET updated_at = $1`, time.Now().Add(-time.Minute*16))
	tx.Commit()
	defer cleanupDB()

	goodCookie := GoodCookier{Value: authToken}
	_, err := authenticate(goodCookie)
	if err == nil {
		t.Error("Expected error not to be nil")
	}
	errString := "sql: no rows in result set"
	if err.Error() != errString {
		t.Errorf("Expected error to equal '%v', but was '%v'\n", errString, err)
	}
}

func TestCreateUser(t *testing.T) {
	// success: new user in database
	// failure: username already exists
}

func TestLogin(t *testing.T) {
	// success: auth token in cookie
	// failure: username incorrect
	// failure: password incorrect
	// failure: internal error (many options...)
}

func TestLogout(t *testing.T) {
	// success: auth token is now invalid
	// failure: invalid auth token in cookie
	// failure: internal error (couldn't invalidate authentications)
}
