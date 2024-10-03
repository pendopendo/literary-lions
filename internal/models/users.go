package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	//"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver import
	"golang.org/x/crypto/bcrypt"
)

// Define a new User struct. Notice how the field names and types align
// with the columns in the database "users" table?
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Define a new UserModel struct which wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the "users" table.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
             VALUES (?, ?, ?, datetime('now'))`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))

	if err != nil {
		// if errors.Is(err, sqlite3.ErrConstraintUnique) { 
		// 	return   ErrDuplicateEmail // ErrDuplicateEmail should be a custom error defined for this case
		// } else {
			return    err
		// }
	}



	// if err != nil {
	// 	// Check for unique constraint violations, etc.
	// 	fmt.Println("sign up error:",  err)
	// 	return err
	// }

	  return nil
}

// Get retrieves a user based on their ID from the "users" table.
func (m *UserModel) Get(id int) (User, error) {
	// Prepare the SQL statement to select the user based on id.
	stmt := `SELECT id, name, email, hashed_password, created FROM users WHERE id = ?`

	// Create a User instance to hold the retrieved data.
	var user User

	// Execute the query and scan the result into the User instance.
	err := m.DB.QueryRow(stmt, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Created,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord // ErrNoRecord should be a custom error defined for this case
		} else {
			return User{}, err
		}
	}

	// If everything is fine, return the User instance.
	return user, nil
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {

	// Retrieve the id and hashed password associated with the given email. If
	// no matching email exists we return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Otherwise, the password is correct. Return the user ID.
	return id, nil

}

// We'll use the Exists method to check if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

func (m *UserModel) InsertUser(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// Update the SQL statement to use datetime('now') for SQLite
	stmt := `INSERT INTO users (name, email, hashed_password, created)
             VALUES(?, ?, ?, datetime('now'))`

	// Execute the statement
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// Check for SQLite constraint violations
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			// Check if the error message relates to the unique email constraint
			if strings.Contains(err.Error(), "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}
