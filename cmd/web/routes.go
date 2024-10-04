package main

import (
	"net/http"
)

// Define a type for middleware functions.
type Middleware func(http.Handler) http.Handler

// Function to chain multiple middleware functions.
func chainMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	// Apply the middleware in reverse order.
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Define your middleware functions as slices.
	dynamicMiddlewares := []Middleware{
		app.sessionManager.LoadAndSave,
	}

	protectedMiddlewares := []Middleware{
		app.sessionManager.LoadAndSave,
		app.requireAuthentication,
	}

	standardMiddlewares := []Middleware{
		commonHeaders,
	}

	// Use chainMiddleware to apply middleware to handlers.
	mux.Handle("GET /{$}", chainMiddleware(http.HandlerFunc(app.home), dynamicMiddlewares...))
	mux.Handle("GET /post/view/{id}", chainMiddleware(http.HandlerFunc(app.postView), dynamicMiddlewares...))

	mux.Handle("GET /post/create/{category_id}", chainMiddleware(http.HandlerFunc(app.postCreate), protectedMiddlewares...))
	mux.Handle("POST /post/create/{category_id}", chainMiddleware(http.HandlerFunc(app.postCreatePost), protectedMiddlewares...))

	mux.Handle("POST /comment/create/{post_id}", chainMiddleware(http.HandlerFunc(app.commentCreatePost), protectedMiddlewares...))
	mux.Handle("GET /category/view/{id}", chainMiddleware(http.HandlerFunc(app.categoryView), dynamicMiddlewares...))

	// User routes
	mux.Handle("GET /user/signup", chainMiddleware(http.HandlerFunc(app.userSignup), dynamicMiddlewares...))
	mux.Handle("POST /user/signup", chainMiddleware(http.HandlerFunc(app.userSignupPost), dynamicMiddlewares...))
	mux.Handle("GET /user/login", chainMiddleware(http.HandlerFunc(app.userLogin), dynamicMiddlewares...))
	mux.Handle("POST /user/login", chainMiddleware(http.HandlerFunc(app.userLoginPost), dynamicMiddlewares...))
	mux.Handle("POST /user/logout", chainMiddleware(http.HandlerFunc(app.userLogoutPost), protectedMiddlewares...))

	// Reaction routes
	mux.Handle("POST /reaction/create/{post_id}", chainMiddleware(http.HandlerFunc(app.reactionCreatePostPost), protectedMiddlewares...))
	mux.Handle("POST /reaction/create/{post_id}/{comment_id}", chainMiddleware(http.HandlerFunc(app.reactionCreateCommentPost), protectedMiddlewares...))

	// Catch-all route for 404 errors
	mux.HandleFunc("/", app.notFound)

	// Wrap the entire mux with the standard middlewares. as as
	return chainMiddleware(mux, standardMiddlewares...)
}
