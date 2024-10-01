package handlers

import (
	"html/template"
	"literary-lions/database"
	"literary-lions/models"
	"net/http"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	//if authenticated siis jah, kui mitte mine logini
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("pages/create_post.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryID, _ := strconv.Atoi(r.FormValue("category_id"))
		userID := getUserIDFromSession(r) // Assume this function gets user ID from the session

		err := models.CreatePost(database.DB, title, content, userID, categoryID)
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func ViewPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts(database.DB)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	tmpl.Execute(w, posts)
}

func FilterPostsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("category")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	posts, err := models.GetPostsByCategory(database.DB, categoryID)
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	tmpl.Execute(w, posts)
}

func FilterPostsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(r)
	posts, err := models.GetPostsByUser(database.DB, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user posts", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	tmpl.Execute(w, posts)
}

func FilterLikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromSession(r)
	posts, err := models.GetLikedPosts(database.DB, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve liked posts", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	tmpl.Execute(w, posts)
}
