package handlers

import (
	"html/template"
	"net/http"
	//"log"
)

/*func serverError(w http.ResponseWriter, err error) {
	log.Printf("Server error: %v", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}*/

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
