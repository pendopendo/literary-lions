package handlers

import (
    "net/http"
    "html/template"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/404.html"))
    w.WriteHeader(http.StatusNotFound)
    tmpl.Execute(w, nil)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/500.html"))
    w.WriteHeader(http.StatusInternalServerError)
    tmpl.Execute(w, nil)
}
