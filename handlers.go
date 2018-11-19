package main

import (
	"bytes"
	"fmt"
	"html"
	"log"
	"net/http"
	"text/template"
)

func ShowRootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
}

func ShowRegistrationHandler(w http.ResponseWriter, r *http.Request) {
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
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Printf("%s %q: %+v\n", r.Method, r.URL.Path, r.Form)

	newUser := User{
		Username: r.Form["username"][0],
		Password: r.Form["password"][0],
	}

	stmt, err := db.PrepareNamed(`INSERT INTO  users (username, password) VALUES (:username, :password)`)
	if err != nil {
		http.Error(w, "Error registering", http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(newUser)
	if err != nil {
		http.Error(w, "Error registering", http.StatusInternalServerError)
		return
	}

	err = db.Get(&newUser, "SELECT * FROM users WHERE username=$1", newUser.Username)
	if err != nil {
		http.Error(w, "Error finding newly registered user", http.StatusInternalServerError)
		return
	}

	log.Printf("User added: %+v\n", newUser)

	fmt.Fprintln(w, "Success!")
}
