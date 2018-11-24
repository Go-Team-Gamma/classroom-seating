package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

func ShowRoot(w http.ResponseWriter, r *http.Request) {
	logHandlerIntro(r.Method, r.URL.Path, r.Form)
	renderPage(w, r, "templates/index.tmpl", "Home")
}

func ShowRegistration(w http.ResponseWriter, r *http.Request) {
	logHandlerIntro(r.Method, r.URL.Path, r.Form)
	renderPage(w, r, "templates/registration.tmpl", "Registration")
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	logHandlerIntro(r.Method, r.URL.Path, r.Form)
	renderPage(w, r, "templates/login.tmpl", "Login")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	logHandlerIntro(r.Method, r.URL.Path, r.Form)

	newUser := User{
		Username: r.Form["username"][0],
		Password: r.Form["password"][0],
	}

	row := db.QueryRowx(
		`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`,
		newUser.Username,
		newUser.Password,
	)
	err := row.Scan(&newUser.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error registering", http.StatusInternalServerError)
		return
	}

	log.Printf("User added: %+v\n", newUser)
	fmt.Fprintln(w, "Success!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	logHandlerIntro(r.Method, r.URL.Path, r.Form)

	user := User{
		Username: r.Form["username"][0],
		Password: r.Form["password"][0],
	}

	// Find the user.
	err := db.Get(
		&user,
		`SELECT * FROM users WHERE username = $1 AND password = $2`,
		user.Username,
		user.Password,
	)
	if err != nil || !user.Id.Valid {
		log.Println(err)
		http.Error(w, "Error logging in", http.StatusBadRequest)
		return
	}

	// Already checked the case where this fails above.
	userId, _ := user.Id.Value()
	authToken, err := findLoginAuthToken(userId.(string))

	if err == nil {
		// We have an auth token that is still valid, let's use that.
		_, err = db.Exec(
			`UPDATE authentications SET updated_at = $1 WHERE token = $2 AND deleted_at IS NULL`,
			time.Now(),
			authToken,
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	} else {
		// Create a new authentication for this user.
		row := db.QueryRowx(`INSERT INTO authentications (user_id) VALUES ($1) RETURNING token`, userId)
		err = row.Scan(&authToken)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error logging in", http.StatusInternalServerError)
			return
		}
	}

	// Invalidate all outstanding authentications for this user.
	_, err = db.Exec(
		`UPDATE authentications SET deleted_at = $1 WHERE user_id = $2 AND deleted_at IS NULL and token != $3`,
		time.Now(),
		userId,
		authToken,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    "authtoken",
		Value:   authToken,
		Expires: time.Now().UTC().Add(time.Minute * 15),
	}
	http.SetCookie(w, cookie)

	log.Printf("Logged in as: %+v\n", user)
	fmt.Fprintln(w, "Success!")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	logHandlerIntro(r.Method, r.URL.Path, r.Form)
	fmt.Fprintln(w, "Success!")

	userId, err := authenticate(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Couldn't authenticate", http.StatusBadRequest)
	}

	_, err = db.Exec(
		`UPDATE authentications SET deleted_at = $1 WHERE user_id = $2 AND deleted_at IS NULL`,
		time.Now(),
		userId,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Couldn't clear authentications", http.StatusInternalServerError)
		return
	}
}

//-- Private, internal helpers.

func logHandlerIntro(requestMethod, requestPath string, requestData url.Values) {
	log.Printf("%s %q: %+v\n", requestMethod, requestPath, requestData)
}

func authenticate(w http.ResponseWriter, r *http.Request) (string, error) {
	authCookie, err := r.Cookie("authtoken")
	if err != nil {
		return "", err
	}

	var userId string
	err = db.Get(
		&userId,
		`SELECT user_id FROM authentications WHERE token = $1 AND deleted_at IS NULL AND updated_at > $2`,
		authCookie.Value,
		time.Now().Add(-time.Minute*15),
	)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Printf("Authenticated as: %v\n", userId)
	return userId, nil
}

func findLoginAuthToken(userId string) (string, error) {
	// Get the last authentication.
	var authToken string
	err := db.Get(
		&authToken,
		`SELECT token FROM authentications WHERE user_id = $1 AND deleted_at IS NULL AND updated_at > $2`,
		userId,
		time.Now().Add(-time.Minute*15),
	)
	if err != nil {
		return "", err
	}

	log.Printf("findLoginAuthToken: authToken: %s\n", authToken)
	return authToken, nil
}

func renderPage(w http.ResponseWriter, r *http.Request, templateName, title string) {
	userId, _ := authenticate(w, r)

	t := template.Must(
		template.ParseFiles(
			"templates/header.tmpl",
			"templates/navigation.tmpl",
			"templates/footer.tmpl",
			templateName))

	data := PageData{
		Title:         title,
		Authenticated: userId != "",
	}

	err := t.ExecuteTemplate(w, strings.ToLower(title), data)

	if err != nil {
		log.Println(err)
	}
}
