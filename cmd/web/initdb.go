package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func InitSQLDB(db *sql.DB) {
	// Create the category table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS category (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the comment table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comment (
			id INTEGER PRIMARY KEY,
			user INTEGER NOT NULL,
			post INTEGER NOT NULL,
			text TEXT NOT NULL,
			created TEXT NOT NULL,
			FOREIGN KEY(user) REFERENCES users(id),
			FOREIGN KEY(post) REFERENCES post(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the post table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS post (
			id INTEGER PRIMARY KEY,
			user INTEGER NOT NULL,
			title TEXT NOT NULL,
			text TEXT NOT NULL,
			category INTEGER NOT NULL,
			created TEXT NOT NULL,
			FOREIGN KEY(user) REFERENCES users(id),
			FOREIGN KEY(category) REFERENCES category(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			hashed_password CHAR(60) NOT NULL,
			created DATETIME NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Adding a UNIQUE constraint to the email and username columns
	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS users_uc_email ON users (email);
		CREATE UNIQUE INDEX IF NOT EXISTS users_uc_name ON users (name);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the reactions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY,
			reaction TEXT NOT NULL,
			comment INTEGER,
			post INTEGER,
			user INTEGER NOT NULL,
			FOREIGN KEY(comment) REFERENCES comment(id),
			FOREIGN KEY(post) REFERENCES post(id),
			FOREIGN KEY(user) REFERENCES users(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the session table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			token CHAR(43) PRIMARY KEY,
			data BLOB NOT NULL,
			expiry DATETIME NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create an index on the expiry column
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert predefined categories into the category table
	_, err = db.Exec(`
		INSERT INTO category (id, name) VALUES
		(1, 'Fiction'),
		(2, 'Non-Fiction'),
		(3, 'Science Fiction'),
		(4, 'Fantasy')
		ON CONFLICT(id) DO NOTHING;
	`)
	if err != nil {
		log.Fatal(err)
	}
}
