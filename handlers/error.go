package handlers

import (
	"html/template"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("pages/404.html"))
	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, nil)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("pages/500.html"))
	w.WriteHeader(http.StatusInternalServerError)
	tmpl.Execute(w, nil)
}
