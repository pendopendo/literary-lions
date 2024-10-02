package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave)
	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request dreaches the file server.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /post/view/{id}", dynamic.ThenFunc(app.postView)) // postituse vaatamiseks
	// http://localhost:4000/post/view/1
	mux.Handle("GET /post/create/{category_id}", protected.ThenFunc(app.postCreate))      //postituse loomine
	mux.Handle("POST /post/create/{category_id}", protected.ThenFunc(app.postCreatePost)) //postituse loomine ID
	//comment
	//mux.Handle("GET /comment/create/{post_id}", app.commentCreate)      //kommi loomine
	mux.Handle("POST /comment/create/{post_id}", protected.ThenFunc(app.commentCreatePost)) //kommi loomine
	mux.Handle("GET /category/view/{id}", dynamic.ThenFunc(app.categoryView))               // postituse vaatamiseks. Kui vajutab link või läheb otse siis see

	//user

	// Add the five new routes, all of which use our 'dynamic' middleware chain.
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	//reaktsioon LiKE
	mux.Handle("POST /reaction/create/{post_id}", protected.ThenFunc(app.reactionCreatePostPost)) // postituse vaatamiseks. Kui vajutab link või läheb otse siis see
	mux.Handle("POST /reaction/create/{post_id}/{comment_id}", protected.ThenFunc(app.reactionCreateCommentPost))
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(commonHeaders)

	// Pass the servemux as the 'next' parameter to the commonHeaders middleware.
	// Because commonHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
