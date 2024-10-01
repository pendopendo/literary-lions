package main

import (
	"database/sql"
	"fmt"
	"literary-lions/database"
	"literary-lions/handlers"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error
	database.DB, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(database.DB)
	defer database.DB.Close()

	database.InitSQLDB(database.DB)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.MainHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.Handle("/create-post",handlers.CreatePostHandler)
	http.Handle("/logout",(handlers.LogoutHandler))
	http.HandleFunc("/filter/category", handlers.FilterPostsByCategoryHandler)
	http.HandleFunc("/filter/user", handlers.FilterPostsByUserHandler)
	http.HandleFunc("/filter/liked", handlers.FilterLikedPostsHandler)
	http.HandleFunc("/404", handlers.NotFoundHandler)
	http.HandleFunc("/500", handlers.InternalServerErrorHandler)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
