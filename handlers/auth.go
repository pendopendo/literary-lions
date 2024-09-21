package handlers

import (
	"database/sql"
	"fmt"
	"literary-lions-forum/models"
	"net/http"
	"text/template"
	"time"
)

type Post struct {
	Title    string
	Category string
	Content  string
}

type MainPage struct {
	Posts      []Post
	Categories []string
}

var db *sql.DB

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		hashedPassword, err := models.HashPassword(password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Pass the db connection to the CreateUser function
		err = models.CreateUser(db, username, email, hashedPassword)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := models.AuthenticateUser(db, email, password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		sessionToken, err := models.CreateSession(db, user.ID)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	posts := []Post{
		{
			Title:    "Programming for dummies",
			Category: "Science",
			Content:  "türa",
		},
		{
			Title:    "Miljon Miksi",
			Category: "Science",
			Content:  "Kuu on tekkinud maakera ja marsi suuruse taevakehade kokkupõrkest väljapaisatud sademete tagajärjel",
		},
		{
			Title:    "We're cooked",
			Category: "Pop Culture",
			Content:  "Deez nuts",
		},
	}
	categories := []string{"Crime", "Science", "Drama"}

	if r.Method == "GET" {
		// Create a new template
		tmpl, err := template.ParseFiles("templates/homepage.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, MainPage{Posts: posts, Categories: categories})
		if err != nil {
			fmt.Println(err)
		}
	}
}
