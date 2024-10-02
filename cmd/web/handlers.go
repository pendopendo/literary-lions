package main

import (
	"errors"
	"fmt"
	"html/template" // New import
	"net/http"
	"strconv"

	//"log"           // New import

	"snippetbox.alexedwards.net/internal/models"
)

//handler võtab requesti sisse, teeb sellega midagi ja annab vastuse (osad ei tee req midagi)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//(29.09.24) seda jrg rida pole enam vaja
	//w.Header().Add("Server", "Go") //response headersis näha browseris (tühikudi ei tohi olla)

	//-----------   VANA   -----------------//
	// kategooriad kuvamien ilma html
	categories, err := app.categories.GetRow()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	//----------- VANA  lõpp --------------//

	// for _, category := range categories {
	// 	fmt.Fprintf(w, "home see katid:%+v\n", category) //home lehele hektel need
	// }

	// //postitused korraks
	// posts, err := app.posts.Get()
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// for _, post := range posts {
	// 	fmt.Fprintf(w, "\nhome see POSTID:%+v\n", post) //home lehele hektel need
	// }

	//fmt.Println("SEE: ", posts)

	// Create an instance of a templateData struct holding the snippet data.
	//seda hakkab siis html kasutama range kaudu
	//saame tänu sellele anda nii post kui comment
	data := templateData{
		Categories:      categories,
		IsAuthenticated: app.isAuthenticated(r),
	}

	// Initialize a slice containing the paths to the view.tmpl file,
	// plus the base layout and navigation partial that we made earlier.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// Parse the template files... loeb failid sisse
	ts, err := template.ParseFiles(files...) // "..." see on 3 argumenti! Iga stringi rida see on nagu files[0], files[1], files[2] täppe kolm alati, see tähenda lihtsalt, et ta loeb jöärjest argumendid sise
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//fmt.Println("Data:", data)
	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter?
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

}

func (app *application) postView(w http.ResponseWriter, r *http.Request) { //http.Request - see tuleb sisse
	id, err := strconv.Atoi(r.PathValue("id")) //konverdib strigni numbriks. URL sees olev ID
	//id - nr ja err vea staatus
	if err != nil || id < 1 {
		/*
			fmt.Fprintf(w, "kas tuleb %d...\n", id) //endale
			fmt.Fprintf(w, "%s", err)//endale
		*/
		http.NotFound(w, r) //käsk http.NotFound võtab selle 404 page not found
		return
	}

	//fmt.Fprintf(w, "Display a specific post with ID %d...", id)

	post, err := app.posts.GetRow(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	//MUDELIGA SUHTLUS
	comments, err := app.comments.GetCommentsForPost(id) // app.comments - sID 240928_184230 sealt
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	//fmt.Println("kommid:", comments)
	// Fetch likes and dislikes for this post
	likes, dislikes, err := app.reactions.CountPostLikesDislikes(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Create an instance of a templateData struct holding the snippet data.
	//seda hakkab siis html kasutama range kaudu
	//saame tänu sellele anda nii post kui comment
	data := templateData{
		Post:            post,
		Comments:        comments, //Comments on template.go failis pandud nimeks sID 240928_214152
		Likes:           likes,    // Pass likes count to the template
		Dislikes:        dislikes, // Pass dislikes count to the template
		IsAuthenticated: app.isAuthenticated(r),
	}

	// Initialize a slice containing the paths to the view.tmpl file,
	// plus the base layout and navigation partial that we made earlier.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	// Parse the template files... loeb failid sisse
	ts, err := template.ParseFiles(files...) // "..." see on 3 argumenti! Iga stringi rida see on nagu files[0], files[1], files[2] täppe kolm alati, see tähenda lihtsalt, et ta loeb jöärjest argumendid sise
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//fmt.Println("Data:", data)
	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter?
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

}

func (app *application) postCreate(w http.ResponseWriter, r *http.Request) {
	//(30.09.24) See on ainult formi näitamine

	categoryID, err := strconv.Atoi(r.PathValue("category_id")) //konverdib strigni numbriks. URL sees olev ID
	//id - nr ja err vea staatus
	if err != nil || categoryID < 1 {
		/*
			fmt.Fprintf(w, "kas tuleb %d...\n", id) //endale
			fmt.Fprintf(w, "%s", err)//endale
		*/
		http.NotFound(w, r) //käsk http.NotFound võtab selle 404 page not found
		return
	}

	//saame tänu sellele anda nii post kui comment
	data := templateData{
		CategoryID:      categoryID,
		IsAuthenticated: app.isAuthenticated(r),
	}

	//w.Write([]byte("postCreate"))
	// Initialize a slice containing the paths to the view.tmpl file,
	// plus the base layout and navigation partial that we made earlier.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/create.tmpl", //see
	}

	// Parse the template files... loeb failid sisse
	ts, err := template.ParseFiles(files...) // "..." see on 3 argumenti! Iga stringi rida see on nagu files[0], files[1], files[2]
	//täppe kolm alati, see tähenda lihtsalt, et ta loeb jöärjest argumendid sise
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//fmt.Println("Data:", data)
	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter?
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

}

// postituse loomine kirjutamine serverisse (30.09.24)
func (app *application) postCreatePost(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("postCreatePost"))
	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError() helper to
	// send a 400 Bad Request response to the user.

	categoryID, err := strconv.Atoi(r.PathValue("category_id")) //konverdib strigni numbriks. URL sees olev ID
	//id - nr ja err vea staatus
	if err != nil || categoryID < 1 {
		/*
			fmt.Fprintf(w, "kas tuleb %d...\n", id) //endale
			fmt.Fprintf(w, "%s", err)//endale
		*/
		http.NotFound(w, r) //käsk http.NotFound võtab selle 404 page not found
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	//leoeme mis SISESTATI
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually covert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
	// Request response if the conversion fails.

	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil {
	//     app.clientError(w, http.StatusBadRequest)
	//     return
	// }

	id, err := app.posts.Insert(categoryID, title, content, userID) //, expires
	//saame siit tagasi ID postituse oma mille lõime
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther) //teeb postituse lahti mille just lõid
}

// func (app *application) commentCreate(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("commentCreate"))
// }

func (app *application) commentCreatePost(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("commentCreate"))
	//ise
	//vaja saada ju posti ID millega see siduda
	//url seest on see id nt /2
	postID, err := strconv.Atoi(r.PathValue("post_id")) //konverdib strigni numbriks. URL sees olev ID
	// id - nr ja err vea staatus
	if err != nil || postID < 1 {
		/*
			fmt.Fprintf(w, "kas tuleb %d...\n", id) //endale
			fmt.Fprintf(w, "%s", err)//endale
		*/
		http.NotFound(w, r) //käsk http.NotFound võtab selle 404 page not found
		return
	}

	//see
	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	// leoeme mis SISESTATI
	content := r.PostForm.Get("content") //<textarea name='content'></textarea>  view.tmpl - see on see content
	fmt.Println("sisu kommil", content)

	// Kasutaja info pärimine
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.comments.InsertComment(postID, content, userID) //, expires
	// saame siit tagasi ID postituse oma mille lõime
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther) //teeb postituse lahti mille just lõid

}

func (app *application) categoryView(w http.ResponseWriter, r *http.Request) {
	//testime(27.09.24)
	//ühe kategooria postituste vaatamine (29.09.2024)

	id, err := strconv.Atoi(r.PathValue("id")) //loeb id //id numbriks
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the categoryModel's Get() method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	/*
		category, err := app.categories.GetRow()
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				app.serverError(w, r, err)
			}
			return
		}

		// Write the category data as a plain-text HTTP response body.
		fmt.Fprintf(w, "KATS: %+v", category)

	*/

	//SIIN PROOVIN  KOKKU PANNA

	//fmt.Fprintf(w, "Display a specific post with ID %d...", id)

	post, err := app.posts.GetPostsForCategory(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// //MUDELIGA SUHTLUS
	// comments, err := app.comments.GetCommentsForPost(id) // app.comments - sID 240928_184230 sealt
	// if err != nil {
	// 	if errors.Is(err, models.ErrNoRecord) {
	// 		http.NotFound(w, r)
	// 	} else {
	// 		app.serverError(w, r, err)
	// 	}
	// 	return
	// }

	// Assuming 'post' is a slice containing the posts
	// for _, p := range post {
	// 	fmt.Fprintf(w, "Postitus: %v\n", p)
	// }

	// Create an instance of a templateData struct holding the snippet data.
	//seda hakkab siis html kasutama range kaudu
	//saame tänu sellele anda nii post kui comment
	data := templateData{
		CategoryID:      id,
		Posts:           post,
		IsAuthenticated: app.isAuthenticated(r),
	}

	// Initialize a slice containing the paths to the view.tmpl file,
	// plus the base layout and navigation partial that we made earlier.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/categoryView.tmpl",
	}

	// Parse the template files... loeb failid sisse
	ts, err := template.ParseFiles(files...) // "..." see on 3 argumenti! Iga stringi rida see on nagu files[0], files[1], files[2] täppe kolm alati, see tähenda lihtsalt, et ta loeb jöärjest argumendid sise
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//fmt.Println("Data:", data)
	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter?
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

}

//----------------------------------------//
//                 USERS                  //
//----------------------------------------//
//(30.09.24)

// Create a new userSignupForm struct.
type userSignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	//validator.Validator `form:"-"`
	FieldErrors map[string]string // This map will hold any validation error messages
}

// Update the handler so it displays the signup page.
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	// Initialize an empty templateData struct
	data := templateData{
		Form: userSignupForm{
			FieldErrors: make(map[string]string), // Initialize the FieldErrors map
		},
		IsAuthenticated: app.isAuthenticated(r),
	}

	app.render(w, r, http.StatusOK, "./ui/html/pages/signup.tmpl", data)
	fmt.Println("user lehele mindud")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Extract form values
	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// Call the Insert method to save the user details
	err = app.users.Insert(name, email, password)
	if err != nil {
		// Handle potential duplicate email error or any other error
		if errors.Is(err, models.ErrDuplicateEmail) {
			// You may want to re-render the signup page with an error message
			http.Error(w, "This email address is already registered", http.StatusConflict)
			return
		}
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to a welcome page or login page after successful signup
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {

	//see on selle nupu vajutamine
	//fmt.Fprintln(w, "Display a form for logging in a user...")
	data := templateData{
		Form: userSignupForm{
			FieldErrors: make(map[string]string), // Initialize the FieldErrors map
		},
		IsAuthenticated: app.isAuthenticated(r),
	}
	fmt.Println("vajutasid login  nuppu")

	app.render(w, r, http.StatusOK, "./ui/html/pages/login.tmpl", data)

}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Authenticate and login the user...")
	// Decode the form data into the userLoginForm struct.

	// Parse the form data from the request
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Extract form values

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {

			app.userLogin(w, r) //kutsume selle siit ülalt

		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	fmt.Println("vajutasid Login rohelist nuppu")
	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther) //home

}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Logout the user...")

	// Use the RenewToken() method on the current session to change the session
	// ID again.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	// Redirect the user to the application home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// LIKE THIS LIKE
func (app *application) reactionCreatePostPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("reactionCreatePost")

	//w.Write([]byte("commentCreate"))
	//ise
	//vaja saada ju posti ID millega see siduda
	//url seest on see id nt /2
	postID, err := strconv.Atoi(r.PathValue("post_id")) //konverdib strigni numbriks. URL sees olev ID
	fmt.Println("reactionCreatePost postID ", postID)
	// id - nr ja err vea staatus
	if err != nil || postID < 1 {
		/*
			fmt.Fprintf(w, "kas tuleb %d...\n", id) //endale
			fmt.Fprintf(w, "%s", err)//endale
		*/
		fmt.Println("Error1")
		http.NotFound(w, r) //käsk http.NotFound võtab selle 404 page not found
		return
	}

	//see
	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error2")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	// leoeme mis SISESTATI
	reaction := r.PostForm.Get("reaction") //<textarea name='content'></textarea>  view.tmpl - see on see content
	fmt.Println("reaktsioon", reaction)

	// Kasutaja info pärimine
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	_, err = app.reactions.InsertReaction(reaction, postID, userID) //, expires
	// saame siit tagasi ID postituse oma mille lõime
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther) //teeb postituse lahti mille just lõid

}

func (app *application) reactionCreateCommentPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("reactionCreateCommentPost")

	//w.Write([]byte("commentCreate"))
	//ise
	//vaja saada ju posti ID millega see siduda
	//url seest on see id nt /2
	commentID, err := strconv.Atoi(r.PathValue("comment_id")) //konverdib strigni numbriks. URL sees olev ID
	fmt.Println("reactionCreateCommentPost commentID ", commentID)
	// id - nr ja err vea staatus
	if err != nil || commentID < 1 {
		/*
			fmt.Fprintf(w, "kas tuleb %d...\n", id) //endale
			fmt.Fprintf(w, "%s", err)//endale
		*/
		fmt.Println("Error1")
		http.NotFound(w, r) //käsk http.NotFound võtab selle 404 page not found
		return
	}

	postID, err := strconv.Atoi(r.PathValue("post_id")) //konverdib strigni numbriks. URL sees olev ID
	fmt.Println("reactionCreateCommentPost postID ", postID)
	// id - nr ja err vea staatus
	if err != nil || postID < 1 {
		/*
			fmt.Fprintf(w, "kas tuleb %d...\n", id) //endale
			fmt.Fprintf(w, "%s", err)//endale
		*/
		fmt.Println("Error1")
		http.NotFound(w, r) //käsk http.NotFound võtab selle 404 page not found
		return
	}

	//see
	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error2")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	// leoeme mis SISESTATI
	reaction := r.PostForm.Get("reaction") //<textarea name='content'></textarea>  view.tmpl - see on see content
	fmt.Println("reaktsioon", reaction)

	// Kasutaja info pärimine
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	_, err = app.reactions.InsertReactionComment(reaction, commentID, userID) //, expires
	// saame siit tagasi ID postituse oma mille lõime
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther) //teeb postituse lahti mille just lõid

}
