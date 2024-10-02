package main

import (
	"bytes"
	"html/template"
	"net/http"
)

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		page, // Use the 'page' parameter here
	}

	// Parse the template files
	ts, err := template.ParseFiles(files...) // "..." expands to files[0], files[1], files[2]
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	// Execute the template named "base" with the provided data
	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Set the status code
	w.WriteHeader(status)

	// Write the buffer to the response writer
	buf.WriteTo(w)
}

// Return true if the current request is from an authenticated user, otherwise
// return false.
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		//CurrentYear: time.Now().Year(),
		//Flash:       app.sessionManager.PopString(r.Context(), "flash"),
		// Add the authentication status to the template data.
		IsAuthenticated: app.isAuthenticated(r),
	}
}
