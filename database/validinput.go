package database
/*
import (
	"net/http"
)

//Error handling function

//checks for valid input while registering a user

func checkForValidInput(w http.ResponseWriter, username, password, email string) {
	 errors := make(map[string]string) // Store validation errors.

	// Check minimal password length.
	if len(password) < 6 {
		errors["password"] = "Please enter password with at least 6 characters"
	}

	// Validate email format.
	if !isValidEmail(email) {
		errors["email"] = "Please enter a correct email"
	}

	// Check if email or username fields are empty.
	if email == "" {
		errors["email"] = "Email field required"
	}
	if username == "" {
		errors["username"] = "Username field required"
	}
	if password == "" {
		errors["password"] = "Password field required"
	}

	// Check if the username already exists in the database.
	stmtForUser, err := db.Prepare("SELECT username FROM users WHERE username = ?;")
	if err != nil {
		serverError(w, err)
		return nil
	}
	defer stmtForUser.Close()

	var userExists string
	err = stmtForUser.QueryRow(username).Scan(&userExists)
	if err != nil && err != sql.ErrNoRows {
		serverError(w, err)
		return nil
	}

	// Check if the email already exists in the database.
	stmtForEmail, err := db.Prepare("SELECT email FROM users WHERE email = ?;")
	if err != nil {
		serverError(w, err)
		return nil
	}
	defer stmtForEmail.Close()

	var emailExists string
	err = stmtForEmail.QueryRow(email).Scan(&emailExists)
	if err != nil && err != sql.ErrNoRows {
		serverError(w, err)
		return nil
	}

	// Add error messages if username or email already exist.
	if userExists != "" {
		errors["username"] = "Username already taken, "
	}
	if emailExists != "" {
		errors["email"] = "Email is already taken"
	}

	return errors
}*/
