package database

import (
	"database/sql"
	"log"
)

var DB *sql.DB

// generates the tables in the database
func InitSQLDB(db *sql.DB) {
	// Create the users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			bio TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the sessions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_token TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the categories table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the posts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			category_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the comments table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the likes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER,
			comment_id INTEGER,
			like BOOLEAN,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (comment_id) REFERENCES comments(id),
			CONSTRAINT unique_like UNIQUE (user_id, post_id, comment_id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the posts_categories table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts_categories (
			post_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (category_id) REFERENCES categories(id),
			CONSTRAINT unique_post_category UNIQUE (post_id, category_id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the posts_likes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts_likes (
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			like BOOLEAN,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			CONSTRAINT unique_post_like UNIQUE (post_id, user_id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the comments_likes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments_likes (
			comment_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			like BOOLEAN,
			FOREIGN KEY (comment_id) REFERENCES comments(id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			CONSTRAINT unique_comment_like UNIQUE (comment_id, user_id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Create the posts_comments table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts_comments (
			post_id INTEGER NOT NULL,
			comment_id INTEGER NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (comment_id) REFERENCES comments(id),
			CONSTRAINT unique_post_comment UNIQUE (post_id, comment_id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
