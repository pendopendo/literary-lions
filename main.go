package main

import (
	"database/sql"
	"literary-lions/database"
	"literary-lions/handlers"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	database.InitSQLDB(db)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.MainHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	//http.Handle("/create-post", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePostHandler)))
	//http.Handle("/logout", middleware.AuthMiddleware(http.HandlerFunc(handlers.LogoutHandler)))
	http.HandleFunc("/filter/category", handlers.FilterPostsByCategoryHandler)
	http.HandleFunc("/filter/user", handlers.FilterPostsByUserHandler)
	http.HandleFunc("/filter/liked", handlers.FilterLikedPostsHandler)
	http.HandleFunc("/404", handlers.NotFoundHandler)
	http.HandleFunc("/500", handlers.InternalServerErrorHandler)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
